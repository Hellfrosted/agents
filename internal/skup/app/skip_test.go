package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExecute_savesSkip_whenSkipCommandFindsUpdate(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	stateDir := filepath.Join(root, "state")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfile())
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "old\n")
	runner := &appFakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": appTarArchive(t, appTarFile{name: "skills/alpha/SKILL.md", contents: "new\n"}),
		},
		hashes: map[string]string{"skills/alpha": "hash-alpha"},
	}

	// When
	code := Execute(context.Background(), Request{
		Argv0:     "sk-up",
		Args:      []string{"-s", "alpha", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", stateDir},
		Env:       map[string]string{"HOME": "/home/alice"},
		Stdout:    &stdout,
		Stderr:    &stderr,
		GitRunner: runner,
		Now:       func() time.Time { return time.Date(2026, 6, 20, 8, 30, 0, 0, time.UTC) },
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
	if !strings.Contains(stdout.String(), "SKIPPED alpha") {
		t.Fatalf("stdout missing skip line: %q", stdout.String())
	}
	got := readAppSkips(t, stateDir)
	entry := got.Skips["alpha"]
	if entry.RemoteHash != "hash-alpha" || entry.SourceURL != "https://github.com/example/skills.git" || entry.SkippedAt != "2026-06-20T08:30:00Z" {
		t.Fatalf("skip entry = %#v", entry)
	}
}

func TestExecute_doesNotSaveSkip_whenSkillIsCurrent(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	stateDir := filepath.Join(root, "state")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfile())
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "same\r\n")
	runner := &appFakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": appTarArchive(t, appTarFile{name: "skills/alpha/SKILL.md", contents: "same\n"}),
		},
		hashes: map[string]string{"skills/alpha": "hash-alpha"},
	}

	// When
	code := Execute(context.Background(), Request{
		Argv0:     "sk-up",
		Args:      []string{"-s", "alpha", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", stateDir},
		Env:       map[string]string{"HOME": "/home/alice"},
		Stdout:    &stdout,
		Stderr:    &bytes.Buffer{},
		GitRunner: runner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "OK      alpha has no current update to skip") {
		t.Fatalf("stdout missing current line: %q", stdout.String())
	}
	got := readAppSkips(t, stateDir)
	if len(got.Skips) != 0 {
		t.Fatalf("skips = %#v, want empty", got.Skips)
	}
}

func TestExecute_removesSkip_whenUnskipCommandProvided(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	root := t.TempDir()
	stateDir := filepath.Join(root, "state")
	writeAppFile(t, filepath.Join(stateDir, "skips.json"), `{"skips":{"alpha":{"remoteHash":"hash-alpha","sourceUrl":"https://github.com/example/skills.git","skippedAt":"2026-06-20T08:30:00Z"}}}`)

	// When
	code := Execute(context.Background(), Request{
		Argv0:  "sk-up",
		Args:   []string{"-u", "alpha", "--agents-home", filepath.Join(root, "agents"), "--cache-dir", filepath.Join(root, "cache"), "--state-dir", stateDir},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &bytes.Buffer{},
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "OK      removed saved skip for alpha") {
		t.Fatalf("stdout missing unskip line: %q", stdout.String())
	}
	got := readAppSkips(t, stateDir)
	if len(got.Skips) != 0 {
		t.Fatalf("skips = %#v, want empty", got.Skips)
	}
}

func TestExecute_writesSkipsJSON_whenSkipsCommandRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	root := t.TempDir()
	stateDir := filepath.Join(root, "state")
	writeAppFile(t, filepath.Join(stateDir, "skips.json"), `{"skips":{"alpha":{"remoteHash":"hash-alpha","sourceUrl":"https://github.com/example/skills.git","skippedAt":"2026-06-20T08:30:00Z"}}}`)

	// When
	code := Execute(context.Background(), Request{
		Argv0:  "sk-up",
		Args:   []string{"-S", "--json", "--agents-home", filepath.Join(root, "agents"), "--cache-dir", filepath.Join(root, "cache"), "--state-dir", stateDir},
		Env:    map[string]string{"HOME": "/home/alice"},
		Stdout: &stdout,
		Stderr: &bytes.Buffer{},
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var got struct {
		Statuses []struct {
			Name       string `json:"name"`
			Status     string `json:"status"`
			RemoteHash string `json:"remoteHash"`
		} `json:"statuses"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, stdout.String())
	}
	if len(got.Statuses) != 1 || got.Statuses[0].Name != "alpha" || got.Statuses[0].Status != "skipped" || got.Statuses[0].RemoteHash != "hash-alpha" {
		t.Fatalf("statuses = %#v", got.Statuses)
	}
}

type appSkipsDocument struct {
	Skips map[string]struct {
		RemoteHash string `json:"remoteHash"`
		SourceURL  string `json:"sourceUrl"`
		SkippedAt  string `json:"skippedAt"`
	} `json:"skips"`
}

func readAppSkips(t *testing.T, stateDir string) appSkipsDocument {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(stateDir, "skips.json"))
	if os.IsNotExist(err) {
		return appSkipsDocument{Skips: map[string]struct {
			RemoteHash string `json:"remoteHash"`
			SourceURL  string `json:"sourceUrl"`
			SkippedAt  string `json:"skippedAt"`
		}{}}
	}
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	var got appSkipsDocument
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v", err)
	}
	return got
}
