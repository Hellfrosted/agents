package plan

import (
	"fmt"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
)

func DryRun(parsed cli.Parsed, resolved config.Resolved) (output.Summary, error) {
	actions, err := actionsFor(parsed, resolved)
	if err != nil {
		return output.Summary{}, err
	}
	return output.Summary{
		OK:         true,
		Command:    string(parsed.Command),
		Entrypoint: string(parsed.Entrypoint),
		DryRun:     true,
		Actions:    actions,
	}, nil
}

func actionsFor(parsed cli.Parsed, resolved config.Resolved) ([]output.PlannedAction, error) {
	switch parsed.Command {
	case cli.CommandInstallSource:
		return installSourceActions(parsed.Targets), nil
	case cli.CommandInstall:
		return installActions(parsed.Targets), nil
	case cli.CommandRemove:
		return removeActions(parsed.Targets, resolved), nil
	case cli.CommandSkip:
		return namedActions("write-skip", parsed.Targets), nil
	case cli.CommandUnskip:
		return namedActions("remove-skip", parsed.Targets), nil
	default:
		return nil, fmt.Errorf("plan dry run for %s: unsupported command", parsed.Command)
	}
}

func installSourceActions(targets []string) []output.PlannedAction {
	actions := make([]output.PlannedAction, 0, len(targets))
	for _, target := range targets {
		actions = append(actions, output.PlannedAction{Action: "install-source", Target: target})
	}
	return actions
}

func installActions(targets []string) []output.PlannedAction {
	if len(targets) == 0 {
		return []output.PlannedAction{
			{Action: "compare-upstream", Target: "changed-or-missing-unskipped"},
			{Action: "install-updates", Target: "changed-or-missing-unskipped"},
		}
	}
	actions := make([]output.PlannedAction, 0, len(targets))
	for _, name := range targets {
		actions = append(actions, output.PlannedAction{Action: "install", Name: name})
	}
	return actions
}

func removeActions(targets []string, resolved config.Resolved) []output.PlannedAction {
	actions := make([]output.PlannedAction, 0, len(targets)*4)
	for _, name := range targets {
		actions = append(actions,
			output.PlannedAction{Action: "remove", Name: name},
			output.PlannedAction{Action: "remove-directory", Name: name, Path: filepath.Join(resolved.AgentsHome, "skills", name)},
			output.PlannedAction{Action: "remove-lock-entry", Name: name, Path: filepath.Join(resolved.AgentsHome, ".skill-lock.json")},
			output.PlannedAction{Action: "remove-skip", Name: name, Path: filepath.Join(resolved.StateDir, "skips.json")},
		)
	}
	return actions
}

func namedActions(action string, targets []string) []output.PlannedAction {
	actions := make([]output.PlannedAction, 0, len(targets))
	for _, name := range targets {
		actions = append(actions, output.PlannedAction{Action: action, Name: name})
	}
	return actions
}
