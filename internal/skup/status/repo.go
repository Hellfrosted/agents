package status

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/output"
)

func ensureRepo(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource, progress ProgressFunc) error {
	if _, err := os.Stat(filepath.Join(source.repoDir, ".git")); err == nil {
		emitProgress(progress, repoProgressEvent(output.EventFetch, source))
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "remote", "set-url", "origin", source.sourceURL); err != nil {
			return err
		}
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "fetch", "--quiet", "--depth", "1", "--filter=blob:none", "origin", "HEAD"); err != nil {
			return err
		}
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "reset", "--quiet", "--soft", "FETCH_HEAD"); err != nil {
			emitProgress(progress, repoProgressEvent(output.EventRepair, source))
			staleDir, err := moveStaleRepo(source.repoDir)
			if err != nil {
				return err
			}
			emitProgress(progress, repoProgressEvent(output.EventClone, source))
			if err := cloneRepo(ctx, runner, gitPath, source); err != nil {
				return restoreStaleRepo(staleDir, source.repoDir, err)
			}
			if err := sparseCheckoutRepo(ctx, runner, gitPath, source); err != nil {
				return restoreStaleRepo(staleDir, source.repoDir, err)
			}
			_ = os.RemoveAll(staleDir)
			return nil
		}
	} else {
		emitProgress(progress, repoProgressEvent(output.EventClone, source))
		if err := cloneRepo(ctx, runner, gitPath, source); err != nil {
			return err
		}
	}
	return sparseCheckoutRepo(ctx, runner, gitPath, source)
}

func moveStaleRepo(repoDir string) (string, error) {
	staleDir := fmt.Sprintf("%s.stale-%d", repoDir, os.Getpid())
	_ = os.RemoveAll(staleDir)
	if err := os.Rename(repoDir, staleDir); err != nil {
		return "", fmt.Errorf("move stale repo cache: %w", err)
	}
	return staleDir, nil
}

func cloneRepo(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource) error {
	if err := os.MkdirAll(filepath.Dir(source.repoDir), 0o700); err != nil {
		return fmt.Errorf("create repo cache directory: %w", err)
	}
	return runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-c", "core.symlinks=false", "clone", "-c", "core.symlinks=false", "--quiet", "--depth", "1", "--filter=blob:none", "--sparse", source.sourceURL, source.repoDir)
}

func sparseCheckoutRepo(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource) error {
	return runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-c", "core.symlinks=false", "-C", source.repoDir, "sparse-checkout", "set", source.remoteDir)
}

func restoreStaleRepo(staleDir string, repoDir string, cause error) error {
	removeErr := os.RemoveAll(repoDir)
	restoreErr := os.Rename(staleDir, repoDir)
	if removeErr != nil || restoreErr != nil {
		return errors.Join(cause, removeErr, restoreErr)
	}
	return cause
}

func remoteHash(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource) (string, error) {
	result, err := runner.Run(ctx, compare.Command{
		Name: gitPath,
		Args: []string{"-c", "core.autocrlf=false", "-C", source.repoDir, "rev-parse", remoteHashRef(source.remoteDir)},
	})
	if err != nil {
		return "", fmt.Errorf("git rev-parse %s: %w", source.name, err)
	}
	return strings.TrimSpace(string(result.Stdout)), nil
}

func remoteHashRef(remoteDir string) string {
	if remoteDir == "" || remoteDir == "." {
		return "HEAD^{tree}"
	}
	return "HEAD:" + remoteDir
}

func runGit(ctx context.Context, runner compare.CommandRunner, gitPath string, args ...string) error {
	if _, err := runner.Run(ctx, compare.Command{Name: gitPath, Args: args}); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}
