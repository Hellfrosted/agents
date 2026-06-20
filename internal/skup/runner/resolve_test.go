package runner

import (
	"errors"
	"reflect"
	"testing"
)

func TestResolveSkillsCommand_usesOverride_whenProvided(t *testing.T) {
	// Given
	input := ResolveInput{
		Override: `"/opt/skills runner" add`,
		Lookup:   lookupNone,
	}

	// When
	got, err := ResolveSkillsCommand(input)

	// Then
	if err != nil {
		t.Fatalf("ResolveSkillsCommand returned error: %v", err)
	}
	want := Command{Program: "/opt/skills runner", Args: []string{"add"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command = %#v, want %#v", got, want)
	}
}

func TestResolveSkillsCommand_prefersFallbackOrder_whenOverrideMissing(t *testing.T) {
	// Given
	input := ResolveInput{
		Lookup: mapLookup(map[string]string{
			"bunx": "/usr/bin/bunx",
			"npx":  "/usr/bin/npx",
		}),
	}

	// When
	got, err := ResolveSkillsCommand(input)

	// Then
	if err != nil {
		t.Fatalf("ResolveSkillsCommand returned error: %v", err)
	}
	want := Command{Program: "/usr/bin/bunx", Args: []string{"skills@latest"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("command = %#v, want %#v", got, want)
	}
}

func TestResolveSkillsCommand_returnsNoRunnerError_whenNoFallbackExists(t *testing.T) {
	// Given
	input := ResolveInput{Lookup: lookupNone}

	// When
	_, err := ResolveSkillsCommand(input)

	// Then
	if !errors.Is(err, ErrNoSkillsRunner) {
		t.Fatalf("error = %v, want ErrNoSkillsRunner", err)
	}
}

func TestResolveSkillsCommand_rejectsUnclosedQuote_whenOverrideMalformed(t *testing.T) {
	// Given
	input := ResolveInput{
		Override: `"pnpm dlx skills@latest`,
		Lookup:   lookupNone,
	}

	// When
	_, err := ResolveSkillsCommand(input)

	// Then
	if !errors.Is(err, ErrInvalidCommand) {
		t.Fatalf("error = %v, want ErrInvalidCommand", err)
	}
}

func lookupNone(string) (string, bool) {
	return "", false
}

func mapLookup(paths map[string]string) LookupFunc {
	return func(name string) (string, bool) {
		path, ok := paths[name]
		return path, ok
	}
}
