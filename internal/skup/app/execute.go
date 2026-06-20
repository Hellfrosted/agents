package app

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/inventory"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/plan"
)

const (
	ExitOK      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

type Request struct {
	Argv0  string
	Args   []string
	Env    map[string]string
	Stdout io.Writer
	Stderr io.Writer

	GitRunner  compare.CommandRunner
	ToolRunner compare.CommandRunner
	Now        func() time.Time
}

func Execute(ctx context.Context, request Request) int {
	if err := ctx.Err(); err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}

	parsed, err := cli.Parse(cli.Input{Argv0: request.Argv0, Args: request.Args})
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}

	if parsed.Command == cli.CommandHelp {
		writeText(request.Stdout, cli.HelpText(parsed.Entrypoint))
		return ExitOK
	}

	resolved, err := config.Resolve(config.ResolveInput{
		Options:  configOptions(parsed),
		Env:      request.Env,
		Platform: config.Platform{},
	})
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}

	if parsed.DryRun {
		return executeDryRun(parsed, resolved, request)
	}
	if parsed.Command == cli.CommandList {
		return executeList(parsed, resolved, request)
	}
	if parsed.Command == cli.CommandStatus {
		return executeStatus(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandDiff {
		return executeDiff(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandOpenDiff {
		return executeOpenDiff(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandInstall || parsed.Command == cli.CommandInstallSource {
		return executeInstall(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandRemove {
		return executeRemove(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandSkip {
		return executeSkip(ctx, parsed, resolved, request)
	}
	if parsed.Command == cli.CommandUnskip {
		return executeUnskip(parsed, resolved, request)
	}
	if parsed.Command == cli.CommandSkips {
		return executeSkips(parsed, resolved, request)
	}

	writeError(request.Stderr, fmt.Errorf("command %q not implemented yet", parsed.Command))
	return ExitUsage
}

func executeList(parsed cli.Parsed, resolved config.Resolved, request Request) int {
	names, err := inventory.ListInstalled(resolved.AgentsHome)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}

	switch parsed.Output {
	case cli.OutputJSON:
		summary := output.Summary{
			OK:         true,
			Command:    string(parsed.Command),
			Entrypoint: string(parsed.Entrypoint),
			Statuses:   skillStatuses(names),
		}
		if err := output.WriteSummaryJSON(request.Stdout, summary); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	case cli.OutputJSONL:
		for _, name := range names {
			event := output.Event{Type: output.EventTypeStatus, Name: name, Status: output.StatusOK}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
	default:
		for _, name := range names {
			writeText(request.Stdout, name+"\n")
		}
	}
	return ExitOK
}

func executeDryRun(parsed cli.Parsed, resolved config.Resolved, request Request) int {
	summary, err := plan.DryRun(parsed, resolved)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}

	switch parsed.Output {
	case cli.OutputJSON:
		if err := output.WriteSummaryJSON(request.Stdout, summary); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	case cli.OutputJSONL:
		event := output.Event{Type: output.EventTypeSummary, OK: summary.OK}
		if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	default:
		for _, action := range summary.Actions {
			writeText(request.Stdout, fmt.Sprintf("PLAN    %s\n", action.Action))
		}
	}
	return ExitOK
}

func skillStatuses(names []string) []output.SkillStatus {
	statuses := make([]output.SkillStatus, 0, len(names))
	for _, name := range names {
		statuses = append(statuses, output.SkillStatus{Name: name, Status: output.StatusOK})
	}
	return statuses
}

func configOptions(parsed cli.Parsed) config.Options {
	return config.Options{
		AgentsHome:    parsed.AgentsHome,
		CacheDir:      parsed.CacheDir,
		StateDir:      parsed.StateDir,
		SkillsCommand: parsed.SkillsCommand,
		DiffTool:      parsed.DiffTool,
		Color:         parsed.Color,
	}
}

func writeError(writer io.Writer, err error) {
	if writer == nil || err == nil {
		return
	}
	fmt.Fprintln(writer, err)
}

func writeText(writer io.Writer, text string) {
	if writer == nil {
		return
	}
	fmt.Fprint(writer, text)
}
