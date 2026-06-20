package compare

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/output"
)

type SkillInput struct {
	GitPath      string
	Repo         string
	RemoteDir    string
	InstalledDir string
	ExportDir    string
}

func CompareGitSkill(ctx context.Context, runner CommandRunner, input SkillInput) (Result, error) {
	if _, err := os.Stat(input.InstalledDir); err != nil {
		if os.IsNotExist(err) {
			return Result{Status: output.StatusMissing}, nil
		}
		return Result{}, fmt.Errorf("stat installed skill: %w", err)
	}

	if err := resetExportDir(input.ExportDir); err != nil {
		return Result{}, err
	}
	if err := ExportGitTree(ctx, runner, ExportInput{
		GitPath:   input.GitPath,
		Repo:      input.Repo,
		RemoteDir: input.RemoteDir,
		Dest:      input.ExportDir,
	}); err != nil {
		return Result{}, err
	}

	return CompareDirs(input.InstalledDir, exportedCompareDir(input.ExportDir, input.RemoteDir))
}

func resetExportDir(path string) error {
	if path == "" || path == "." {
		return fmt.Errorf("reset export directory: unsafe path %q", path)
	}
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("remove export directory: %w", err)
	}
	if err := os.MkdirAll(path, 0o700); err != nil {
		return fmt.Errorf("create export directory: %w", err)
	}
	return nil
}

func exportedCompareDir(exportDir string, remoteDir string) string {
	if remoteDir == "" || remoteDir == "." {
		return exportDir
	}
	return filepath.Join(exportDir, filepath.FromSlash(remoteDir))
}
