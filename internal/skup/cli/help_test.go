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

func TestHelpText_explainsShortFlags_whenEntrypointIsSkUp(t *testing.T) {
	// Given
	entrypoint := EntrypointShort

	// When
	got := HelpText(entrypoint)

	// Then
	wantContains := []string{
		"Commands:",
		"  -h, --help                 Show this help.",
		"  -l, --list                 List installed skills.",
		"  -g, --global               Check installed skills against upstream.",
		"  -d, --diff <skill>         Print a terminal diff for one skill.",
		"  -z, --zed [skill...]       Open changed skills in the configured diff tool.",
		"  -i, --install [skill...]   Install updates for named skills, or all changed skills.",
		"  -I, --install-source <source...>",
		"  -s, --skip <skill>         Skip the current upstream revision for one skill.",
		"  -u, --unskip <skill>       Clear saved skip state for one skill.",
		"  -S, --skips                List saved skip state.",
		"  -r, --remove <skill...>    Remove installed skills and lockfile entries.",
		"Output and control flags:",
		"  --json                     Write one final JSON summary to stdout.",
		"  --jsonl                    Write newline-delimited JSON events to stdout.",
		"  --dry-run                  Show planned mutations without changing files.",
		"  --agents-home <path>       Override the installed agents home.",
		"  --cache-dir <path>         Override the repository cache directory.",
		"  --state-dir <path>         Override updater state storage.",
		"  --skills-command <cmd>     Override the delegated Skills CLI command.",
		"  --diff-tool <cmd>          Override the GUI diff command.",
		"  --color auto|always|never  Control colored human output.",
		"  --no-color                 Disable colored human output.",
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
