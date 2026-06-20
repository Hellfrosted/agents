package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

func executeStatus(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	result, err := status.Check(ctx, gitRunner(request), status.Input{
		GitPath:    "git",
		AgentsHome: resolved.AgentsHome,
		CacheDir:   resolved.CacheDir,
		StateDir:   resolved.StateDir,
		Targets:    parsed.Targets,
	})
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
			Statuses:   result.Statuses,
		}
		if err := output.WriteSummaryJSON(request.Stdout, summary); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	case cli.OutputJSONL:
		for _, status := range result.Statuses {
			event := output.Event{Type: output.EventTypeStatus, Name: status.Name, Status: status.Status, RemoteHash: status.RemoteHash}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
	default:
		for _, status := range result.Statuses {
			writeText(request.Stdout, fmt.Sprintf("%-7s %s\n", strings.ToUpper(string(status.Status)), status.Name))
		}
	}
	return ExitOK
}

func gitRunner(request Request) compare.CommandRunner {
	if request.GitRunner != nil {
		return request.GitRunner
	}
	return compare.ExecRunner{}
}
