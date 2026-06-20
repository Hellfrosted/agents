package cli

import "strings"

func HelpText(entrypoint Entrypoint) string {
	if entrypoint == EntrypointLong {
		return longHelp()
	}
	return shortHelp()
}

func shortHelp() string {
	lines := []string{
		"Usage:",
		"  sk-up -h",
		"  sk-up -l",
		"  sk-up -g",
		"  sk-up -d <skill>",
		"  sk-up -z [skill...]",
		"  sk-up -i [skill...]",
		"  sk-up -I <source...>",
		"  sk-up -s <skill>",
		"  sk-up -u <skill>",
		"  sk-up -S",
		"  sk-up -r <skill...>",
		"",
		"Common flags:",
		"  --json",
		"  --jsonl",
		"  --dry-run",
		"  --agents-home <path>",
		"  --cache-dir <path>",
		"  --state-dir <path>",
		"  --skills-command <cmd>",
		"  --diff-tool <cmd>",
		"  --color auto|always|never",
		"  --no-color",
	}
	return strings.Join(lines, "\n") + "\n"
}

func longHelp() string {
	lines := []string{
		"Usage:",
		"  skills-updates --help",
		"  skills-updates --list",
		"  skills-updates --global",
		"  skills-updates --diff <skill>",
		"  skills-updates --zed [skill...]",
		"  skills-updates --install [skill...]",
		"  skills-updates --install-source <source...>",
		"  skills-updates --skip <skill>",
		"  skills-updates --unskip <skill>",
		"  skills-updates --skips",
		"  skills-updates --remove <skill...>",
		"",
		"Common flags:",
		"  --json",
		"  --jsonl",
		"  --dry-run",
		"  --agents-home <path>",
		"  --cache-dir <path>",
		"  --state-dir <path>",
		"  --skills-command <cmd>",
		"  --diff-tool <cmd>",
		"  --color auto|always|never",
		"  --no-color",
	}
	return strings.Join(lines, "\n") + "\n"
}
