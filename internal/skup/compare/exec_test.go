package compare

import (
	"context"
	"errors"
	"os"
	"testing"
)

func TestExecRunner_capturesStdout_whenCommandSucceeds(t *testing.T) {
	// Given
	runner := ExecRunner{}
	command := Command{
		Name: os.Args[0],
		Args: []string{"-test.run=TestExecHelperProcess", "--", "ok"},
	}

	// When
	got, err := runner.Run(context.Background(), command)

	// Then
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if string(got.Stdout) != "hello\n" {
		t.Fatalf("Stdout = %q, want hello", got.Stdout)
	}
}

func TestExecRunner_returnsCommandError_whenCommandFails(t *testing.T) {
	// Given
	runner := ExecRunner{}
	command := Command{
		Name: os.Args[0],
		Args: []string{"-test.run=TestExecHelperProcess", "--", "fail"},
	}

	// When
	_, err := runner.Run(context.Background(), command)

	// Then
	var commandErr *CommandError
	if !errors.As(err, &commandErr) {
		t.Fatalf("Run error = %v, want CommandError", err)
	}
	if commandErr.Command.Name != os.Args[0] {
		t.Fatalf("CommandError command = %#v", commandErr.Command)
	}
	if commandErr.Exit != 7 {
		t.Fatalf("CommandError Exit = %d, want 7", commandErr.Exit)
	}
	if string(commandErr.Stdout) != "partial\n" {
		t.Fatalf("CommandError Stdout = %q, want partial", commandErr.Stdout)
	}
}

func TestExecHelperProcess(t *testing.T) {
	index := helperArgIndex(os.Args)
	if index < 0 {
		return
	}
	switch os.Args[index] {
	case "ok":
		if _, err := os.Stdout.WriteString("hello\n"); err != nil {
			os.Exit(2)
		}
		os.Exit(0)
	case "fail":
		if _, err := os.Stdout.WriteString("partial\n"); err != nil {
			os.Exit(2)
		}
		if _, err := os.Stderr.WriteString("nope\n"); err != nil {
			os.Exit(2)
		}
		os.Exit(7)
	}
}

func helperArgIndex(args []string) int {
	for i, arg := range args {
		if arg == "--" && i+1 < len(args) {
			return i + 1
		}
	}
	return -1
}
