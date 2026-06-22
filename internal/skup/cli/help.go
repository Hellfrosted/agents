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
		"Commands:",
		"  -h, --help                 Show this help.",
		"  -l, --list                 List installed skills.",
		"  -g, --global               Check installed skills against upstream.",
		"  -d, --diff <skill>         Print a terminal diff for one skill.",
		"  -z, --zed [skill...]       Open changed skills in the configured diff tool.",
		"  -i, --install [skill...]   Install updates for named skills, or all changed skills.",
		"  -I, --install-source <source...>",
		"                             Install skills from explicit source strings.",
		"  -s, --skip <skill>         Skip the current upstream revision for one skill.",
		"  -u, --unskip <skill>       Clear saved skip state for one skill.",
		"  -S, --skips                List saved skip state.",
		"  -r, --remove <skill...>    Remove installed skills and lockfile entries.",
		"",
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
		"Commands:",
		"  --help                         Show this help.",
		"  --list                         List installed skills.",
		"  --global                       Check installed skills against upstream.",
		"  --diff <skill>                 Print a terminal diff for one skill.",
		"  --zed [skill...]               Open changed skills in the configured diff tool.",
		"  --install [skill...]           Install updates for named skills, or all changed skills.",
		"  --install-source <source...>   Install skills from explicit source strings.",
		"  --skip <skill>                 Skip the current upstream revision for one skill.",
		"  --unskip <skill>               Clear saved skip state for one skill.",
		"  --skips                        List saved skip state.",
		"  --remove <skill...>            Remove installed skills and lockfile entries.",
		"",
		"Output and control flags:",
		"  --json                         Write one final JSON summary to stdout.",
		"  --jsonl                        Write newline-delimited JSON events to stdout.",
		"  --dry-run                      Show planned mutations without changing files.",
		"  --agents-home <path>           Override the installed agents home.",
		"  --cache-dir <path>             Override the repository cache directory.",
		"  --state-dir <path>             Override updater state storage.",
		"  --skills-command <cmd>         Override the delegated Skills CLI command.",
		"  --diff-tool <cmd>              Override the GUI diff command.",
		"  --color auto|always|never      Control colored human output.",
		"  --no-color                     Disable colored human output.",
	}
	return strings.Join(lines, "\n") + "\n"
}
