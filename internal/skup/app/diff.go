package app

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/runner"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

func executeDiff(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	statuses, code := checkForDiff(ctx, parsed, resolved, request)
	if code != ExitOK {
		return code
	}
	skillStatus := statuses[0]
	if parsed.Output != cli.OutputHuman {
		return writeDiffSummary(parsed, request, statuses)
	}
	if skillStatus.Status != output.StatusUpdate {
		writeText(request.Stdout, fmt.Sprintf("%-7s %s\n", strings.ToUpper(string(skillStatus.Status)), skillStatus.Name))
		return ExitOK
	}
	writeText(request.Stdout, fmt.Sprintf("\n===== %s =====\n", skillStatus.Name))
	return runTerminalDiff(ctx, skillStatus, request)
}

func executeOpenDiff(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	statuses, code := checkForDiff(ctx, parsed, resolved, request)
	if code != ExitOK {
		return code
	}
	changed := updateStatuses(statuses)
	if len(changed) > 0 {
		if code := runOpenDiff(ctx, resolved, request, changed); code != ExitOK {
			return code
		}
	}
	if parsed.Output != cli.OutputHuman {
		return writeDiffSummary(parsed, request, statuses)
	}
	if len(changed) == 0 {
		writeText(request.Stdout, openDiffNoopLine(parsed, statuses))
	}
	return ExitOK
}

func checkForDiff(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) ([]output.SkillStatus, int) {
	result, err := status.Check(ctx, gitRunner(request), status.Input{
		GitPath:    "git",
		AgentsHome: resolved.AgentsHome,
		CacheDir:   resolved.CacheDir,
		StateDir:   resolved.StateDir,
		Targets:    parsed.Targets,
	})
	if err != nil {
		writeError(request.Stderr, err)
		return nil, ExitUsage
	}
	if len(result.Statuses) == 0 && len(parsed.Targets) > 0 {
		writeError(request.Stderr, fmt.Errorf("skill %s not found in lockfile", parsed.Targets[0]))
		return nil, ExitUsage
	}
	return result.Statuses, ExitOK
}

func runTerminalDiff(ctx context.Context, skillStatus output.SkillStatus, request Request) int {
	command := compare.Command{
		Name: "git",
		Args: []string{
			"-c", "core.autocrlf=false",
			"diff",
			"--ignore-cr-at-eol",
			"--no-index",
			"--color=auto",
			"--",
			skillStatus.InstalledDir,
			skillStatus.CompareDir,
		},
	}
	result, err := gitRunner(request).Run(ctx, command)
	if err != nil {
		var commandErr *compare.CommandError
		if errors.As(err, &commandErr) && commandErr.Exit == 1 {
			writeText(request.Stdout, string(commandErr.Stdout))
			return ExitOK
		}
		writeError(request.Stderr, fmt.Errorf("%s: diff failed: %w", skillStatus.Name, err))
		return ExitFailure
	}
	writeText(request.Stdout, string(result.Stdout))
	return ExitOK
}

func runOpenDiff(ctx context.Context, resolved config.Resolved, request Request, statuses []output.SkillStatus) int {
	command, err := runner.ParseCommand(resolved.DiffTool)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	args := append([]string(nil), command.Args...)
	for _, status := range statuses {
		args = append(args, "--diff", status.InstalledDir, status.CompareDir)
	}
	if request.ToolRunner == nil {
		request.ToolRunner = compare.ExecRunner{}
	}
	if _, err := request.ToolRunner.Run(ctx, compare.Command{Name: command.Program, Args: args}); err != nil {
		writeError(request.Stderr, fmt.Errorf("open diff tool failed: %w", err))
		return ExitFailure
	}
	return ExitOK
}

func writeDiffSummary(parsed cli.Parsed, request Request, statuses []output.SkillStatus) int {
	switch parsed.Output {
	case cli.OutputJSON:
		return writeSummary(request, parsed, statuses)
	case cli.OutputJSONL:
		for _, status := range statuses {
			event := output.Event{Type: output.EventTypeStatus, Name: status.Name, Status: status.Status, SourceURL: status.SourceURL, RemoteHash: status.RemoteHash}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
	}
	return ExitOK
}

func updateStatuses(statuses []output.SkillStatus) []output.SkillStatus {
	changed := make([]output.SkillStatus, 0, len(statuses))
	for _, status := range statuses {
		if status.Status == output.StatusUpdate {
			changed = append(changed, status)
		}
	}
	return changed
}

func openDiffNoopLine(parsed cli.Parsed, statuses []output.SkillStatus) string {
	if len(parsed.Targets) == 1 && len(statuses) == 1 {
		return fmt.Sprintf("%-7s %s\n", strings.ToUpper(string(statuses[0].Status)), statuses[0].Name)
	}
	if len(parsed.Targets) > 0 {
		return "OK      requested skills are up to date\n"
	}
	return "OK      all skills are up to date\n"
}
