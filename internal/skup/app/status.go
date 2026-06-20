package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

func executeStatus(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	progress := status.ProgressFunc(nil)
	if parsed.Output == cli.OutputHuman {
		progress = humanProgress(request.Stderr, resolved.Color)
	}
	result, err := status.Check(ctx, gitRunner(request), status.Input{
		GitPath:    "git",
		AgentsHome: resolved.AgentsHome,
		CacheDir:   resolved.CacheDir,
		StateDir:   resolved.StateDir,
		Targets:    parsed.Targets,
		Progress:   progress,
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
			writeText(request.Stdout, fmt.Sprintf("%s %s\n", humanStatusLabel(status.Status, resolved.Color, request.Stdout), status.Name))
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

func humanStatusLabel(status output.Status, colorMode string, writer io.Writer) string {
	label := fmt.Sprintf("%-7s", strings.ToUpper(string(status)))
	if !statusColorEnabled(colorMode, writer) {
		return label
	}
	if color, ok := statusColor(status); ok {
		return color + label + "\x1b[0m"
	}
	return label
}

func statusColor(status output.Status) (string, bool) {
	switch status {
	case output.StatusOK:
		return "\x1b[32m", true
	case output.StatusUpdate:
		return "\x1b[33m", true
	case output.StatusMissing, output.StatusError:
		return "\x1b[31m", true
	case output.StatusSkipped:
		return "\x1b[36m", true
	default:
		return "", false
	}
}

func statusColorEnabled(colorMode string, writer io.Writer) bool {
	switch colorMode {
	case "always":
		return true
	case "never":
		return false
	default:
		file, ok := writer.(*os.File)
		if !ok {
			return false
		}
		info, err := file.Stat()
		return err == nil && info.Mode()&os.ModeCharDevice != 0
	}
}
