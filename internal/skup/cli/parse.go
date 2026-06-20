package cli

import (
	"fmt"
	"path/filepath"
)

func Parse(input Input) (Parsed, error) {
	parsed := Parsed{
		Entrypoint: detectEntrypoint(input.Argv0),
		Command:    CommandHelp,
		Output:     OutputHuman,
		Color:      "auto",
	}
	structuredOutputSet := false

	args := append([]string(nil), input.Args...)
	for len(args) > 0 {
		arg := args[0]
		args = args[1:]

		if arg == "--" {
			parsed.Targets = append(parsed.Targets, args...)
			break
		}

		if !isOption(arg) {
			parsed.Targets = append(parsed.Targets, arg)
			continue
		}

		var err error
		parsed, args, structuredOutputSet, err = applyOption(parsed, arg, args, structuredOutputSet)
		if err != nil {
			return Parsed{}, err
		}
	}

	if err := validate(parsed); err != nil {
		return Parsed{}, err
	}
	return parsed, nil
}

func detectEntrypoint(argv0 string) Entrypoint {
	name := filepath.Base(argv0)
	switch name {
	case string(EntrypointLong), "skills-updates.exe", "skills-updates.cmd":
		return EntrypointLong
	default:
		return EntrypointShort
	}
}

func isOption(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

func applyOption(parsed Parsed, option string, args []string, structuredOutputSet bool) (Parsed, []string, bool, error) {
	switch option {
	case "-h", "--help":
		parsed.Command = CommandHelp
	case "-l", "--list":
		parsed.Command = CommandList
	case "-g", "--global":
		parsed.Command = CommandStatus
	case "-d", "--diff":
		parsed.Command = CommandDiff
	case "-z", "--zed", "--gui":
		parsed.Command = CommandOpenDiff
	case "-i", "--install", "--install-all", "install-all":
		parsed.Command = CommandInstall
	case "-I", "--install-source":
		parsed.Command = CommandInstallSource
	case "-s", "--skip":
		parsed.Command = CommandSkip
	case "-u", "--unskip":
		parsed.Command = CommandUnskip
	case "-S", "--skips":
		parsed.Command = CommandSkips
	case "-r", "--remove", "--uninstall":
		parsed.Command = CommandRemove
	case "--json":
		if structuredOutputSet && parsed.Output != OutputJSON {
			return Parsed{}, nil, false, fmt.Errorf("%w: choose only one structured output mode", ErrUsage)
		}
		parsed.Output = OutputJSON
		structuredOutputSet = true
	case "--jsonl":
		if structuredOutputSet && parsed.Output != OutputJSONL {
			return Parsed{}, nil, false, fmt.Errorf("%w: choose only one structured output mode", ErrUsage)
		}
		parsed.Output = OutputJSONL
		structuredOutputSet = true
	case "--dry-run":
		parsed.DryRun = true
	case "--no-color":
		parsed.Color = "never"
	case "--agents-home", "--cache-dir", "--state-dir", "--skills-command", "--diff-tool", "--color":
		value, rest, err := takeValue(option, args)
		if err != nil {
			return Parsed{}, nil, false, err
		}
		return applyValueOption(parsed, option, value), rest, structuredOutputSet, nil
	default:
		return Parsed{}, nil, false, fmt.Errorf("%w: unknown option %s", ErrUsage, option)
	}
	return parsed, args, structuredOutputSet, nil
}

func applyValueOption(parsed Parsed, option string, value string) Parsed {
	switch option {
	case "--agents-home":
		parsed.AgentsHome = value
	case "--cache-dir":
		parsed.CacheDir = value
	case "--state-dir":
		parsed.StateDir = value
	case "--skills-command":
		parsed.SkillsCommand = value
	case "--diff-tool":
		parsed.DiffTool = value
	case "--color":
		parsed.Color = value
	}
	return parsed
}

func takeValue(option string, args []string) (string, []string, error) {
	if len(args) == 0 || isOption(args[0]) {
		return "", nil, fmt.Errorf("%w: %s requires a value", ErrUsage, option)
	}
	return args[0], args[1:], nil
}

func validate(parsed Parsed) error {
	if parsed.Output == OutputJSONL && parsed.Command == CommandHelp {
		return fmt.Errorf("%w: --jsonl cannot be used with help", ErrUsage)
	}
	if parsed.Command == CommandDiff && len(parsed.Targets) != 1 {
		return fmt.Errorf("%w: diff requires exactly one skill", ErrUsage)
	}
	if parsed.Command == CommandSkip && len(parsed.Targets) != 1 {
		return fmt.Errorf("%w: skip requires exactly one skill", ErrUsage)
	}
	if parsed.Command == CommandUnskip && len(parsed.Targets) != 1 {
		return fmt.Errorf("%w: unskip requires exactly one skill", ErrUsage)
	}
	if parsed.Command == CommandInstallSource && len(parsed.Targets) == 0 {
		return fmt.Errorf("%w: install-source requires at least one source", ErrUsage)
	}
	if parsed.Command == CommandRemove && len(parsed.Targets) == 0 {
		return fmt.Errorf("%w: remove requires at least one skill", ErrUsage)
	}
	return nil
}
