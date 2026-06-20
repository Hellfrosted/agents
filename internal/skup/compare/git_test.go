package compare

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestExportGitTree_extractsArchiveFromGitStdout(t *testing.T) {
	// Given
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "skills/demo/SKILL.md", contents: "demo\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	runner := &fakeRunner{result: CommandResult{Stdout: archive.Bytes()}}
	dest := t.TempDir()
	input := ExportInput{
		GitPath:   "/usr/bin/git",
		Repo:      "/cache/repo",
		RemoteDir: "skills/demo",
		Dest:      dest,
	}

	// When
	err := ExportGitTree(context.Background(), runner, input)

	// Then
	if err != nil {
		t.Fatalf("ExportGitTree returned error: %v", err)
	}
	wantCommand := Command{
		Name: "/usr/bin/git",
		Args: []string{"-c", "core.autocrlf=false", "-C", "/cache/repo", "archive",
			"--format=tar", "HEAD", "skills/demo"},
	}
	if !reflect.DeepEqual(runner.commands[0], wantCommand) {
		t.Fatalf("command = %#v, want %#v", runner.commands[0], wantCommand)
	}
	raw, err := os.ReadFile(filepath.Join(dest, "skills", "demo", "SKILL.md"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(raw) != "demo\n" {
		t.Fatalf("file contents = %q", raw)
	}
}

func TestExportGitTree_omitsRemoteDir_whenRemoteDirIsDot(t *testing.T) {
	// Given
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	runner := &fakeRunner{result: CommandResult{Stdout: archive.Bytes()}}
	input := ExportInput{GitPath: "git", Repo: "/repo", RemoteDir: ".", Dest: t.TempDir()}

	// When
	err := ExportGitTree(context.Background(), runner, input)

	// Then
	if err != nil {
		t.Fatalf("ExportGitTree returned error: %v", err)
	}
	wantArgs := []string{"-c", "core.autocrlf=false", "-C", "/repo", "archive", "--format=tar", "HEAD"}
	if !reflect.DeepEqual(runner.commands[0].Args, wantArgs) {
		t.Fatalf("args = %#v, want %#v", runner.commands[0].Args, wantArgs)
	}
}

func TestExportGitTree_returnsRunnerError_whenGitFails(t *testing.T) {
	// Given
	runner := &fakeRunner{err: errors.New("git failed")}
	input := ExportInput{GitPath: "git", Repo: "/repo", RemoteDir: ".", Dest: t.TempDir()}

	// When
	err := ExportGitTree(context.Background(), runner, input)

	// Then
	if err == nil || !errors.Is(err, runner.err) {
		t.Fatalf("ExportGitTree error = %v, want runner error", err)
	}
}

type fakeRunner struct {
	commands []Command
	result   CommandResult
	err      error
}

func (r *fakeRunner) Run(_ context.Context, command Command) (CommandResult, error) {
	r.commands = append(r.commands, command)
	return r.result, r.err
}
