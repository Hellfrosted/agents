package app

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/compare"
)

func TestExecute_removesSkillAndCleansLocalState(t *testing.T) {
	// Given
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	stateDir := filepath.Join(root, "state")
	lockPath := filepath.Join(agentsHome, ".skill-lock.json")
	writeAppFile(t, lockPath, `{"version":1,"workspace":{"owner":"local"},"skills":{"alpha":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/alpha/SKILL.md"},"beta":{"sourceUrl":"https://github.com/example/skills.git","skillPath":"skills/beta/SKILL.md"}}}`)
	writeAppFile(t, filepath.Join(agentsHome, "skills", "alpha", "SKILL.md"), "alpha\n")
	writeAppFile(t, filepath.Join(stateDir, "skips.json"), `{"skips":{"alpha":{"remoteHash":"hash-alpha","sourceUrl":"https://github.com/example/skills.git","skippedAt":"2026-06-20T08:30:00Z"}}}`)
	toolRunner := &appFakeGitRunner{}

	// When
	code := Execute(context.Background(), Request{
		Argv0:      "sk-up",
		Args:       []string{"-r", "alpha", "--skills-command", "skills", "--agents-home", agentsHome, "--cache-dir", filepath.Join(root, "cache"), "--state-dir", stateDir},
		Env:        map[string]string{"HOME": "/home/alice"},
		Stdout:     &bytes.Buffer{},
		Stderr:     &bytes.Buffer{},
		ToolRunner: toolRunner,
	})

	// Then
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	wantCommand := compare.Command{
		Name: "skills",
		Args: []string{"remove", "-g", "-y", "--agent", "universal", "--skill", "alpha"},
	}
	if !reflect.DeepEqual(toolRunner.commands[0], wantCommand) {
		t.Fatalf("command = %#v, want %#v", toolRunner.commands[0], wantCommand)
	}
	if _, err := os.Stat(filepath.Join(agentsHome, "skills", "alpha")); !os.IsNotExist(err) {
		t.Fatalf("alpha directory still exists or stat failed differently: %v", err)
	}
	raw, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("ReadFile lockfile returned error: %v", err)
	}
	var lock struct {
		Skills map[string]json.RawMessage `json:"skills"`
	}
	if err := json.Unmarshal(raw, &lock); err != nil {
		t.Fatalf("Unmarshal lockfile returned error: %v", err)
	}
	if _, ok := lock.Skills["alpha"]; ok {
		t.Fatalf("alpha lock entry remained: %s", raw)
	}
	if _, ok := lock.Skills["beta"]; !ok {
		t.Fatalf("beta lock entry missing: %s", raw)
	}
	skips := readAppSkips(t, stateDir)
	if len(skips.Skips) != 0 {
		t.Fatalf("skips = %#v, want empty", skips.Skips)
	}
}
