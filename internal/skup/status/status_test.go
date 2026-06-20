package status

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestCheck_returnsStatusesFromLockfileBackedComparison(t *testing.T) {
	// Given
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	stateDir := filepath.Join(root, "state")
	writeStatusFile(t, filepath.Join(agentsHome, ".skill-lock.json"), lockfileFixture())
	writeStatusFile(t, filepath.Join(stateDir, "skips.json"), skipsFixture())
	writeStatusFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "same\r\n")
	writeStatusFile(t, filepath.Join(agentsHome, "skills", "skipped", "SKILL.md"), "local\n")
	runner := &fakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha":   tarArchive(t, statusTarFile{name: "skills/alpha/SKILL.md", contents: "same\n"}),
			"skills/skipped": tarArchive(t, statusTarFile{name: "skills/skipped/SKILL.md", contents: "remote\n"}),
		},
		hashes: map[string]string{
			"skills/alpha":   "hash-alpha",
			"skills/beta":    "hash-beta",
			"skills/skipped": "hash-skipped",
		},
	}

	// When
	got, err := Check(context.Background(), runner, Input{
		GitPath:    "git",
		AgentsHome: agentsHome,
		CacheDir:   cacheDir,
		StateDir:   stateDir,
	})

	// Then
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	want := []output.SkillStatus{
		{
			Name:         "alpha",
			Status:       output.StatusOK,
			SourceURL:    "https://github.com/example/skills.git",
			RemoteHash:   "hash-alpha",
			InstalledDir: filepath.Join(agentsHome, "skills", "alpha"),
			CompareDir:   filepath.Join(cacheDir, "exports", "alpha", "skills", "alpha"),
		},
		{
			Name:         "beta",
			Status:       output.StatusMissing,
			SourceURL:    "https://github.com/example/skills.git",
			RemoteHash:   "hash-beta",
			InstalledDir: filepath.Join(agentsHome, "skills", "beta"),
			CompareDir:   filepath.Join(cacheDir, "exports", "beta", "skills", "beta"),
		},
		{
			Name:         "skipped",
			Status:       output.StatusSkipped,
			SourceURL:    "https://github.com/example/skills.git",
			RemoteHash:   "hash-skipped",
			InstalledDir: filepath.Join(agentsHome, "skills", "skipped"),
			CompareDir:   filepath.Join(cacheDir, "exports", "skipped", "skills", "skipped"),
		},
	}
	if !reflect.DeepEqual(got.Statuses, want) {
		t.Fatalf("Statuses = %#v, want %#v", got.Statuses, want)
	}
	if len(runner.commands) == 0 {
		t.Fatal("runner commands empty")
	}
}

func TestCheck_filtersTargets_whenTargetsProvided(t *testing.T) {
	// Given
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	writeStatusFile(t, filepath.Join(agentsHome, ".skill-lock.json"), lockfileFixture())
	writeStatusFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "same\n")
	runner := &fakeGitRunner{
		archives: map[string][]byte{
			"skills/alpha": tarArchive(t, statusTarFile{name: "skills/alpha/SKILL.md", contents: "same\n"}),
		},
		hashes: map[string]string{"skills/alpha": "hash-alpha"},
	}

	// When
	got, err := Check(context.Background(), runner, Input{
		GitPath:    "git",
		AgentsHome: agentsHome,
		CacheDir:   cacheDir,
		StateDir:   filepath.Join(root, "state"),
		Targets:    []string{"alpha"},
	})

	// Then
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if len(got.Statuses) != 1 || got.Statuses[0].Name != "alpha" {
		t.Fatalf("Statuses = %#v", got.Statuses)
	}
}

type fakeGitRunner struct {
	commands []compare.Command
	archives map[string][]byte
	hashes   map[string]string
}

func (r *fakeGitRunner) Run(_ context.Context, command compare.Command) (compare.CommandResult, error) {
	r.commands = append(r.commands, command)
	if containsArg(command.Args, "archive") {
		remoteDir := command.Args[len(command.Args)-1]
		return compare.CommandResult{Stdout: r.archives[remoteDir]}, nil
	}
	if containsArg(command.Args, "rev-parse") {
		arg := command.Args[len(command.Args)-1]
		if arg == "HEAD:." {
			return compare.CommandResult{}, fmt.Errorf("root tree must not use HEAD:.")
		}
		remoteDir := "."
		if arg != "HEAD^{tree}" {
			remoteDir = strings.TrimPrefix(arg, "HEAD:")
		}
		return compare.CommandResult{Stdout: []byte(r.hashes[remoteDir] + "\n")}, nil
	}
	return compare.CommandResult{}, nil
}

func lockfileFixture() string {
	return `{
  "version": 1,
  "skills": {
    "alpha": {"sourceUrl": "https://github.com/example/skills.git", "skillPath": "skills/alpha/SKILL.md"},
    "beta": {"sourceUrl": "https://github.com/example/skills.git", "skillPath": "skills/beta/SKILL.md"},
    "skipped": {"sourceUrl": "https://github.com/example/skills.git", "skillPath": "skills/skipped/SKILL.md"}
  }
}`
}

func skipsFixture() string {
	return `{
  "skips": {
    "skipped": {"remoteHash": "hash-skipped", "sourceUrl": "https://github.com/example/skills.git", "skippedAt": "2026-06-20T00:00:00Z"}
  }
}`
}

type statusTarFile struct {
	name     string
	contents string
}

func tarArchive(t *testing.T, file statusTarFile) []byte {
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

func writeStatusFile(t *testing.T, path string, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

func containsArg(args []string, want string) bool {
	for _, arg := range args {
		if arg == want {
			return true
		}
	}
	return false
}
