package plan

import (
	"reflect"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestDryRun_returnsInstallSourceActions_whenInstallSourceCommandProvided(t *testing.T) {
	// Given
	parsed := cli.Parsed{
		Entrypoint: cli.EntrypointShort,
		Command:    cli.CommandInstallSource,
		Targets:    []string{"owner/repo"},
		DryRun:     true,
	}
	resolved := config.Resolved{AgentsHome: "/agents", CacheDir: "/cache", StateDir: "/state"}

	// When
	got, err := DryRun(parsed, resolved)

	// Then
	if err != nil {
		t.Fatalf("DryRun returned error: %v", err)
	}
	want := []output.PlannedAction{
		{Action: "install-source", Target: "owner/repo"},
	}
	if !reflect.DeepEqual(got.Actions, want) {
		t.Fatalf("Actions = %#v, want %#v", got.Actions, want)
	}
	if !got.OK {
		t.Fatal("OK = false, want true")
	}
}

func TestDryRun_returnsNamedInstallActions_whenInstallCommandHasTargets(t *testing.T) {
	// Given
	parsed := cli.Parsed{
		Entrypoint: cli.EntrypointShort,
		Command:    cli.CommandInstall,
		Targets:    []string{"confidence-loop"},
		DryRun:     true,
	}
	resolved := config.Resolved{AgentsHome: "/agents", CacheDir: "/cache", StateDir: "/state"}

	// When
	got, err := DryRun(parsed, resolved)

	// Then
	if err != nil {
		t.Fatalf("DryRun returned error: %v", err)
	}
	want := []output.PlannedAction{
		{Action: "install", Name: "confidence-loop"},
	}
	if !reflect.DeepEqual(got.Actions, want) {
		t.Fatalf("Actions = %#v, want %#v", got.Actions, want)
	}
}

func TestDryRun_returnsRemoveCleanupActions_whenRemoveCommandProvided(t *testing.T) {
	// Given
	parsed := cli.Parsed{
		Entrypoint: cli.EntrypointShort,
		Command:    cli.CommandRemove,
		Targets:    []string{"old-skill"},
		DryRun:     true,
	}
	resolved := config.Resolved{AgentsHome: "/agents", CacheDir: "/cache", StateDir: "/state"}

	// When
	got, err := DryRun(parsed, resolved)

	// Then
	if err != nil {
		t.Fatalf("DryRun returned error: %v", err)
	}
	want := []output.PlannedAction{
		{Action: "remove", Name: "old-skill"},
		{Action: "remove-directory", Name: "old-skill", Path: "/agents/skills/old-skill"},
		{Action: "remove-lock-entry", Name: "old-skill", Path: "/agents/.skill-lock.json"},
		{Action: "remove-skip", Name: "old-skill", Path: "/state/skips.json"},
	}
	if !reflect.DeepEqual(got.Actions, want) {
		t.Fatalf("Actions = %#v, want %#v", got.Actions, want)
	}
}

func TestDryRun_returnsInstallChangedPlan_whenInstallCommandHasNoTargets(t *testing.T) {
	// Given
	parsed := cli.Parsed{
		Entrypoint: cli.EntrypointShort,
		Command:    cli.CommandInstall,
		DryRun:     true,
	}
	resolved := config.Resolved{AgentsHome: "/agents", CacheDir: "/cache", StateDir: "/state"}

	// When
	got, err := DryRun(parsed, resolved)

	// Then
	if err != nil {
		t.Fatalf("DryRun returned error: %v", err)
	}
	want := []output.PlannedAction{
		{Action: "compare-upstream", Target: "changed-or-missing-unskipped"},
		{Action: "install-updates", Target: "changed-or-missing-unskipped"},
	}
	if !reflect.DeepEqual(got.Actions, want) {
		t.Fatalf("Actions = %#v, want %#v", got.Actions, want)
	}
}
