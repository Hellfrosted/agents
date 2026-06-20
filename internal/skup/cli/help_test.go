package cli

import "testing"

func TestHelpText_returnsShortHelp_whenEntrypointIsSkUp(t *testing.T) {
	// Given
	entrypoint := EntrypointShort

	// When
	got := HelpText(entrypoint)

	// Then
	wantContains := []string{
		"Usage:",
		"  sk-up -g",
		"  sk-up -I <source...>",
		"--json",
		"--jsonl",
		"--dry-run",
	}
	for _, want := range wantContains {
		if !contains(got, want) {
			t.Fatalf("HelpText(%q) missing %q in:\n%s", entrypoint, want, got)
		}
	}
}

func TestHelpText_returnsLongHelp_whenEntrypointIsSkillsUpdates(t *testing.T) {
	// Given
	entrypoint := EntrypointLong

	// When
	got := HelpText(entrypoint)

	// Then
	wantContains := []string{
		"Usage:",
		"  skills-updates --global",
		"  skills-updates --install-source <source...>",
		"--skills-command <cmd>",
		"--diff-tool <cmd>",
	}
	for _, want := range wantContains {
		if !contains(got, want) {
			t.Fatalf("HelpText(%q) missing %q in:\n%s", entrypoint, want, got)
		}
	}
}

func contains(text string, part string) bool {
	return len(part) == 0 || index(text, part) >= 0
}

func index(text string, part string) int {
	for i := 0; i+len(part) <= len(text); i++ {
		if text[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
