package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/compare"
)

func TestExecute_delegatesInstallSourceToSkillsCommand(t *testing.T) {
	// Given
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := t.TempDir()
	toolRunner := &appFakeGitRunner{}

	// When
	code := Execute(context.Background(), Request{
		Argv0:      "sk-up",
		Args:       []string{"-I", "owner/repo", "--skills-command", "pnpm dlx skills@latest", "--agents-home", filepath.Join(root, "agents"), "--cache-dir", filepath.Join(root, "cache"), "--state-dir", filepath.Join(root, "state")},
		Env:        map[string]string{"HOME": "/home/alice"},
		Stdout:     &stdout,
		Stderr:     &stderr,
		ToolRunner: toolRunner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	wantCommand := compare.Command{
		Name: "pnpm",
		Args: []string{"dlx", "skills@latest", "add", "owner/repo", "-g", "-y", "--agent", "universal"},
	}
	if !reflect.DeepEqual(toolRunner.commands[0], wantCommand) {
		t.Fatalf("command = %#v, want %#v", toolRunner.commands[0], wantCommand)
	}
	if !strings.Contains(stdout.String(), "INSTALL owner/repo") {
		t.Fatalf("stdout missing install line: %q", stdout.String())
	}
}

func TestExecute_delegatesNamedInstallToSkillsCommand(t *testing.T) {
	// Given
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	writeAppFile(t, filepath.Join(agentsHome, ".skill-lock.json"), appLockfile())
	toolRunner := &appFakeGitRunner{}

	// When
	code := Execute(context.Background(), Request{
		Argv0:      "sk-up",
		Args:       []string{"-i", "alpha", "--skills-command", "skills", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", filepath.Join(root, "state")},
		Env:        map[string]string{"HOME": "/home/alice"},
		Stdout:     &bytes.Buffer{},
		Stderr:     &stderr,
		ToolRunner: toolRunner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0; stderr=%s", code, stderr.String())
	}
	wantCommand := compare.Command{
		Name: "skills",
		Args: []string{"add", "https://github.com/example/skills.git", "-g", "-y", "--agent", "universal", "--skill", "alpha"},
	}
	if !reflect.DeepEqual(toolRunner.commands[0], wantCommand) {
		t.Fatalf("command = %#v, want %#v", toolRunner.commands[0], wantCommand)
	}
}

func TestExecute_installsChangedSkills_whenInstallHasNoTargets(t *testing.T) {
	// Given
	var stderr bytes.Buffer
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
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
		Args:       []string{"-i", "--skills-command", "skills", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", filepath.Join(root, "state")},
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
	if len(toolRunner.commands) != 1 {
		t.Fatalf("tool command count = %d, want 1", len(toolRunner.commands))
	}
	wantCommand := compare.Command{
		Name: "skills",
		Args: []string{"add", "https://github.com/example/skills.git", "-g", "-y", "--agent", "universal", "--skill", "alpha"},
	}
	if !reflect.DeepEqual(toolRunner.commands[0], wantCommand) {
		t.Fatalf("command = %#v, want %#v", toolRunner.commands[0], wantCommand)
	}
}

func TestExecute_preservesExistingLockfileFieldsAfterDelegatedInstall(t *testing.T) {
	// Given
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	lockPath := filepath.Join(agentsHome, ".skill-lock.json")
	writeAppFile(t, lockPath, `{"version":1,"workspace":{"owner":"local"},"skills":{"alpha":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/alpha/SKILL.md","extra":{"keep":true}}}}`)
	toolRunner := &mutatingInstallRunner{lockPath: lockPath}

	// When
	code := Execute(context.Background(), Request{
		Argv0:      "sk-up",
		Args:       []string{"-i", "alpha", "--skills-command", "skills", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", filepath.Join(root, "state")},
		Env:        map[string]string{"HOME": "/home/alice"},
		Stdout:     &bytes.Buffer{},
		Stderr:     &bytes.Buffer{},
		ToolRunner: toolRunner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	raw, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	var got struct {
		Workspace map[string]string `json:"workspace"`
		Skills    map[string]struct {
			Extra map[string]bool `json:"extra"`
		} `json:"skills"`
	}
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("Unmarshal returned error: %v\n%s", err, raw)
	}
	if got.Workspace["owner"] != "local" {
		t.Fatalf("workspace = %#v", got.Workspace)
	}
	if !got.Skills["alpha"].Extra["keep"] {
		t.Fatalf("alpha extra field not preserved: %s", raw)
	}
}

type mutatingInstallRunner struct {
	lockPath string
}

func (r *mutatingInstallRunner) Run(_ context.Context, _ compare.Command) (compare.CommandResult, error) {
	raw := `{"version":1,"skills":{"alpha":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/alpha/SKILL.md"}}}`
	if err := os.WriteFile(r.lockPath, []byte(raw), 0o600); err != nil {
		return compare.CommandResult{}, err
	}
	return compare.CommandResult{}, nil
}
