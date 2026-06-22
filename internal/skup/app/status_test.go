package app

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestExecute_writesStatusJSON_whenGlobalJSONRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
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
		Args:      []string{"-g", "--json", "--agents-home", agentsHome, "--cache-dir", cacheDir, "--state-dir", stateDir},
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
	if len(got.Statuses) != 1 || got.Statuses[0].Status != "ok" || got.Statuses[0].RemoteHash != "hash-alpha" {
		t.Fatalf("statuses = %#v", got.Statuses)
	}
}

func TestExecute_writesStatusLines_whenGlobalHumanRequested(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfileWithSkills("alpha", "beta"))
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "old\n")
	writeAppFile(t, filepath.Join(agentsHome, "skills", "beta", "SKILL.md"), "same\r\n")
	runner := &appFakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": appTarArchive(t, appTarFile{name: "skills/alpha/SKILL.md", contents: "new\n"}),
			"skills/beta":  appTarArchive(t, appTarFile{name: "skills/beta/SKILL.md", contents: "same\n"}),
		},
		hashes: map[string]string{"skills/alpha": "hash-alpha", "skills/beta": "hash-beta"},
	}

	// When
	code := Execute(context.Background(), Request{
		Argv0:     "sk-up",
		Args:      []string{"-g", "--agents-home", agentsHome, "--cache-dir", cacheDir, "--state-dir", filepath.Join(root, "state")},
		Env:       map[string]string{"HOME": "/home/alice"},
		Stdout:    &stdout,
		Stderr:    &stderr,
		GitRunner: runner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "1/2 UPDATE  alpha") {
		t.Fatalf("stdout missing numbered update line: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "2/2 OK      beta") {
		t.Fatalf("stdout missing numbered ok line: %q", stdout.String())
	}
	if !strings.Contains(stderr.String(), "CLONE   https://github.com/example/skills.git") {
		t.Fatalf("stderr missing clone progress: %q", stderr.String())
	}
	if !strings.Contains(stderr.String(), "CHECK   alpha") {
		t.Fatalf("stderr missing check progress: %q", stderr.String())
	}
}

func TestExecute_colorsUpdateStatus_whenColorAlways(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
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
		Args:      []string{"-g", "--color", "always", "--agents-home", agentsHome, "--cache-dir", cacheDir, "--state-dir", filepath.Join(root, "state")},
		Env:       map[string]string{"HOME": "/home/alice"},
		Stdout:    &stdout,
		Stderr:    &bytes.Buffer{},
		GitRunner: runner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "\x1b[33mUPDATE ") {
		t.Fatalf("stdout missing colored update line: %q", stdout.String())
	}
}

func TestHumanLabels_useDistinctColors_whenCheckingAndSkipped(t *testing.T) {
	// Given
	var writer bytes.Buffer

	// When
	checkLabel := humanProgressLine(output.Event{Event: output.EventCompare, Name: "alpha"}, "always", &writer)
	skippedLabel := humanStatusLabel(output.StatusSkipped, "always", &writer)

	// Then
	if !strings.Contains(checkLabel, "\x1b[36mCHECK") {
		t.Fatalf("check label = %q, want cyan CHECK", checkLabel)
	}
	if strings.Contains(skippedLabel, "\x1b[36mSKIPPED") {
		t.Fatalf("skipped label reused CHECK cyan: %q", skippedLabel)
	}
	if !strings.Contains(skippedLabel, "\x1b[35mSKIPPED") {
		t.Fatalf("skipped label = %q, want magenta SKIPPED", skippedLabel)
	}
}

type appFakeGitRunner struct {
	mu         sync.Mutex
	archives   map[string][]byte
	hashes     map[string]string
	commands   []compare.Command
	diffStdout []byte
}

func (r *appFakeGitRunner) Run(_ context.Context, command compare.Command) (compare.CommandResult, error) {
	r.mu.Lock()
	r.commands = append(r.commands, command)
	r.mu.Unlock()
	if hasAppArg(command.Args, "archive") {
		return compare.CommandResult{Stdout: r.archives[command.Args[len(command.Args)-1]]}, nil
	}
	if hasAppArg(command.Args, "rev-parse") {
		remoteDir := strings.TrimPrefix(command.Args[len(command.Args)-1], "HEAD:")
		return compare.CommandResult{Stdout: []byte(r.hashes[remoteDir] + "\n")}, nil
	}
	if hasAppArg(command.Args, "diff") {
		return compare.CommandResult{Stdout: append([]byte(nil), r.diffStdout...)}, nil
	}
	return compare.CommandResult{}, nil
}

func (r *appFakeGitRunner) lastCommandWithArg(arg string) compare.Command {
	r.mu.Lock()
	defer r.mu.Unlock()
	for index := len(r.commands) - 1; index >= 0; index-- {
		if hasAppArg(r.commands[index].Args, arg) {
			return r.commands[index]
		}
	}
	return compare.Command{}
}

func appLockfile() string {
	return `{"version":1,"skills":{"alpha":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/alpha/SKILL.md"}}}`
}

func appLockfileWithSkills(names ...string) string {
	entries := make([]string, 0, len(names))
	for _, name := range names {
		entries = append(entries, `"`+name+`":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/`+name+`/SKILL.md"}`)
	}
	return `{"version":1,"skills":{` + strings.Join(entries, ",") + `}}`
}

type appTarFile struct {
	name     string
	contents string
}

func appTarArchive(t *testing.T, file appTarFile) []byte {
	t.Helper()
	var buffer bytes.Buffer
	tw := tar.NewWriter(&buffer)
	header := &tar.Header{Name: file.name, Mode: 0o600, Size: int64(len(file.contents))}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("WriteHeader returned error: %v", err)
	}
	if _, err := tw.Write([]byte(file.contents)); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	return buffer.Bytes()
}

func writeAppFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

func hasAppArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}
