package cli

import (
	"errors"
	"reflect"
	"testing"
)

func TestParse_returnsStatusCommand_whenSkUpGlobalFlagProvided(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"-g", "--json"}}

	// When
	got, err := Parse(input)

	// Then
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got.Entrypoint != EntrypointShort {
		t.Fatalf("Entrypoint = %q, want %q", got.Entrypoint, EntrypointShort)
	}
	if got.Command != CommandStatus {
		t.Fatalf("Command = %q, want %q", got.Command, CommandStatus)
	}
	if got.Output != OutputJSON {
		t.Fatalf("Output = %q, want %q", got.Output, OutputJSON)
	}
}

func TestParse_returnsInstallSourceCommand_whenShortInstallSourceFlagProvided(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"-I", "owner/repo", "--dry-run"}}

	// When
	got, err := Parse(input)

	// Then
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got.Command != CommandInstallSource {
		t.Fatalf("Command = %q, want %q", got.Command, CommandInstallSource)
	}
	if !got.DryRun {
		t.Fatal("DryRun = false, want true")
	}
	wantTargets := []string{"owner/repo"}
	if !reflect.DeepEqual(got.Targets, wantTargets) {
		t.Fatalf("Targets = %#v, want %#v", got.Targets, wantTargets)
	}
}

func TestParse_keepsInstallArgumentsAsSkillNames_whenShortInstallFlagProvided(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"-i", "owner/repo"}}

	// When
	got, err := Parse(input)

	// Then
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got.Command != CommandInstall {
		t.Fatalf("Command = %q, want %q", got.Command, CommandInstall)
	}
	wantTargets := []string{"owner/repo"}
	if !reflect.DeepEqual(got.Targets, wantTargets) {
		t.Fatalf("Targets = %#v, want %#v", got.Targets, wantTargets)
	}
}

func TestParse_returnsRemoveCommand_whenLongRemoveFlagProvided(t *testing.T) {
	// Given
	input := Input{Argv0: "skills-updates", Args: []string{"--remove", "old-skill"}}

	// When
	got, err := Parse(input)

	// Then
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if got.Entrypoint != EntrypointLong {
		t.Fatalf("Entrypoint = %q, want %q", got.Entrypoint, EntrypointLong)
	}
	if got.Command != CommandRemove {
		t.Fatalf("Command = %q, want %q", got.Command, CommandRemove)
	}
	wantTargets := []string{"old-skill"}
	if !reflect.DeepEqual(got.Targets, wantTargets) {
		t.Fatalf("Targets = %#v, want %#v", got.Targets, wantTargets)
	}
}

func TestParse_rejectsUnknownOption(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"--wat"}}

	// When
	_, err := Parse(input)

	// Then
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("Parse error = %v, want ErrUsage", err)
	}
}

func TestParse_rejectsDiffWithoutExactlyOneTarget(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"-d", "one", "two"}}

	// When
	_, err := Parse(input)

	// Then
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("Parse error = %v, want ErrUsage", err)
	}
}

func TestParse_rejectsBothStructuredOutputModes(t *testing.T) {
	// Given
	input := Input{Argv0: "sk-up", Args: []string{"-g", "--json", "--jsonl"}}

	// When
	_, err := Parse(input)

	// Then
	if !errors.Is(err, ErrUsage) {
		t.Fatalf("Parse error = %v, want ErrUsage", err)
	}
}
