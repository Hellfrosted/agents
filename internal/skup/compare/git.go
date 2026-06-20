package compare

import (
	"bytes"
	"context"
	"fmt"
)

type Command struct {
	Name string
	Args []string
}

type CommandResult struct {
	Stdout []byte
}

type CommandRunner interface {
	Run(ctx context.Context, command Command) (CommandResult, error)
}

type ExportInput struct {
	GitPath   string
	Repo      string
	RemoteDir string
	Dest      string
}

func ExportGitTree(ctx context.Context, runner CommandRunner, input ExportInput) error {
	result, err := runner.Run(ctx, gitArchiveCommand(input))
	if err != nil {
		return fmt.Errorf("git archive: %w", err)
	}
	if err := ExtractTar(bytes.NewReader(result.Stdout), input.Dest); err != nil {
		return fmt.Errorf("extract git archive: %w", err)
	}
	return nil
}

func gitArchiveCommand(input ExportInput) Command {
	args := []string{
		"-c", "core.autocrlf=false",
		"-C", input.Repo,
		"archive",
		"--format=tar",
		"HEAD",
	}
	if input.RemoteDir != "" && input.RemoteDir != "." {
		args = append(args, input.RemoteDir)
	}
	return Command{Name: input.GitPath, Args: args}
}
