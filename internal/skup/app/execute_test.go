package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecute_writesHelpToStdout_whenHelpRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	request := Request{
		Argv0:  "sk-up",
		Args:   []string{"-h"},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// When
	code := Execute(context.Background(), request)

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "sk-up -I <source...>") {
		t.Fatalf("stdout missing short help:\n%s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestExecute_writesDryRunJSONToStdoutOnly_whenStructuredDryRunRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	request := Request{
		Argv0:  "sk-up",
		Args:   []string{"-I", "owner/repo", "--dry-run", "--json", "--agents-home", "/agents", "--cache-dir", "/cache", "--state-dir", "/state"},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// When
	code := Execute(context.Background(), request)

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	var got struct {
		OK      bool `json:"ok"`
		DryRun  bool `json:"dryRun"`
		Actions []struct {
			Action string `json:"action"`
			Target string `json:"target"`
		} `json:"actions"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout.String())
	}
	if !got.OK || !got.DryRun {
		t.Fatalf("summary flags = ok:%v dryRun:%v", got.OK, got.DryRun)
	}
	if len(got.Actions) != 1 || got.Actions[0].Action != "install-source" {
		t.Fatalf("actions = %#v", got.Actions)
	}
}

func TestExecute_writesDryRunJSONLToStdoutOnly_whenStructuredStreamRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	request := Request{
		Argv0:  "sk-up",
		Args:   []string{"-I", "owner/repo", "--dry-run", "--jsonl", "--agents-home", "/agents", "--cache-dir", "/cache", "--state-dir", "/state"},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// When
	code := Execute(context.Background(), request)

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"type":"summary"`) {
		t.Fatalf("stdout missing summary event: %s", stdout.String())
	}
}

func TestExecute_listsInstalledSkillsToStdout_whenListRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	agentsHome := createAgentsHome(t, []string{"zeta", "alpha"})
	request := Request{
		Argv0:  "sk-up",
		Args:   []string{"-l", "--agents-home", agentsHome, "--cache-dir", "/cache", "--state-dir", "/state"},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// When
	code := Execute(context.Background(), request)

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stdout.String() != "alpha\nzeta\n" {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestExecute_listsInstalledSkillsAsJSON_whenListAndJSONRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	agentsHome := createAgentsHome(t, []string{"alpha"})
	request := Request{
		Argv0:  "sk-up",
		Args:   []string{"-l", "--json", "--agents-home", agentsHome, "--cache-dir", "/cache", "--state-dir", "/state"},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	// When
	code := Execute(context.Background(), request)

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	var got struct {
		Statuses []struct {
			Name   string `json:"name"`
			Status string `json:"status"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout.String())
	}
	if len(got.Statuses) != 1 || got.Statuses[0].Name != "alpha" || got.Statuses[0].Status != "ok" {
		t.Fatalf("statuses = %#v", got.Statuses)
	}
}

func createAgentsHome(t *testing.T, names []string) string {
	t.Helper()
	agentsHome := t.TempDir()
	for _, name := range names {
		if err := os.MkdirAll(filepath.Join(agentsHome, "skills", name), 0o700); err != nil {
			t.Fatalf("MkdirAll returned error: %v", err)
		}
	}
	return agentsHome
}
