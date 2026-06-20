package app

import (
	"bytes"
	"context"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/compare"
)

func TestExecute_writesTerminalDiff_whenDiffFindsUpdate(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfile())
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "old\n")
	runner := &appFakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": appTarArchive(t, appTarFile{name: "skills/alpha/SKILL.md", contents: "new\n"}),
		},
		hashes:     map[string]string{"skills/alpha": "hash-alpha"},
		diffStdout: []byte("diff --git a/SKILL.md b/SKILL.md\n"),
	}

	// When
	code := Execute(context.Background(), Request{
		Argv0:     "sk-up",
		Args:      []string{"-d", "alpha", "--agents-home", agentsHome, "--cache-dir", cacheDir, "--state-dir", filepath.Join(root, "state")},
		Env:       map[string]string{"HOME": "/home/alice"},
		Stdout:    &stdout,
		Stderr:    &stderr,
		GitRunner: runner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "===== alpha =====") || !strings.Contains(stdout.String(), "diff --git") {
		t.Fatalf("stdout missing diff content: %q", stdout.String())
	}
	wantDiffArgs := []string{
		"-c", "core.autocrlf=false",
		"diff",
		"--ignore-cr-at-eol",
		"--no-index",
		"--color=auto",
		"--",
		filepath.Join(agentsHome, "skills", "alpha"),
		filepath.Join(cacheDir, "exports", "alpha", "skills", "alpha"),
	}
	if !reflect.DeepEqual(runner.lastCommandWithArg("diff").Args, wantDiffArgs) {
		t.Fatalf("diff args = %#v, want %#v", runner.lastCommandWithArg("diff").Args, wantDiffArgs)
	}
}

func TestExecute_opensDiffToolForUpdates_whenOpenDiffRequested(t *testing.T) {
	// Given
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfile())
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "old\n")
	gitRunner := &appFakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": appTarArchive(t, appTarFile{name: "skills/alpha/SKILL.md", contents: "new\n"}),
		},
		hashes: map[string]string{"skills/alpha": "hash-alpha"},
	}
	toolRunner := &appFakeGitRunner{}

	// When
	code := Execute(context.Background(), Request{
		Argv0:      "sk-up",
		Args:       []string{"-z", "alpha", "--agents-home", agentsHome, "--cache-dir", cacheDir, "--state-dir", filepath.Join(root, "state")},
		Env:        map[string]string{"HOME": "/home/alice"},
		Stdout:     &bytes.Buffer{},
		Stderr:     &stderr,
		GitRunner:  gitRunner,
		ToolRunner: toolRunner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	wantCommand := compare.Command{
		Name: "zed",
		Args: []string{
			"--diff",
			filepath.Join(agentsHome, "skills", "alpha"),
			filepath.Join(cacheDir, "exports", "alpha", "skills", "alpha"),
		},
	}
	if !reflect.DeepEqual(toolRunner.commands[0], wantCommand) {
		t.Fatalf("tool command = %#v, want %#v", toolRunner.commands[0], wantCommand)
	}
}
