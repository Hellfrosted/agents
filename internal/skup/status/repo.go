package status

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Hellfrosted/agents/internal/skup/compare"
)

func ensureRepo(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource) error {
	if _, err := os.Stat(filepath.Join(source.repoDir, ".git")); err == nil {
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "remote", "set-url", "origin", source.sourceURL); err != nil {
			return err
		}
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "fetch", "--quiet", "--depth", "1", "--filter=blob:none", "origin", "HEAD"); err != nil {
			return err
		}
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "reset", "--quiet", "--hard", "FETCH_HEAD"); err != nil {
			return err
		}
	} else {
		if err := os.MkdirAll(filepath.Dir(source.repoDir), 0o700); err != nil {
			return fmt.Errorf("create repo cache directory: %w", err)
		}
		if err := runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "clone", "--quiet", "--depth", "1", "--filter=blob:none", "--sparse", source.sourceURL, source.repoDir); err != nil {
			return err
		}
	}
	return runGit(ctx, runner, gitPath, "-c", "core.autocrlf=false", "-C", source.repoDir, "sparse-checkout", "set", source.remoteDir)
}

func remoteHash(ctx context.Context, runner compare.CommandRunner, gitPath string, source skillSource) (string, error) {
	result, err := runner.Run(ctx, compare.Command{
		Name: gitPath,
		Args: []string{"-c", "core.autocrlf=false", "-C", source.repoDir, "rev-parse", "HEAD:" + source.remoteDir},
	})
	if err != nil {
		return "", fmt.Errorf("git rev-parse %s: %w", source.name, err)
	}
	return strings.TrimSpace(string(result.Stdout)), nil
}

func runGit(ctx context.Context, runner compare.CommandRunner, gitPath string, args ...string) error {
	if _, err := runner.Run(ctx, compare.Command{Name: gitPath, Args: args}); err != nil {
		return fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return nil
}
