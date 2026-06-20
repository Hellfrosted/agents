package compare

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type ExecRunner struct{}

type CommandError struct {
	Command Command
	Stdout  []byte
	Stderr  []byte
	Exit    int
	Err     error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("run %s: %v", e.Command.Name, e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}

func (ExecRunner) Run(ctx context.Context, command Command) (CommandResult, error) {
	cmd := exec.CommandContext(ctx, command.Name, command.Args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		exitCode := -1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		return CommandResult{}, &CommandError{
			Command: command,
			Stdout:  append([]byte(nil), stdout.Bytes()...),
			Stderr:  append([]byte(nil), stderr.Bytes()...),
			Exit:    exitCode,
			Err:     err,
		}
	}
	return CommandResult{Stdout: append([]byte(nil), stdout.Bytes()...)}, nil
}
