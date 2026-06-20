package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/lockfile"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/runner"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

type installItem struct {
	name   string
	source string
}

func executeInstall(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	items, err := installItems(ctx, parsed, resolved, request)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	if len(items) == 0 {
		return writeInstallNoop(parsed, request)
	}
	if err := runInstallItems(ctx, resolved, request, items); err != nil {
		writeError(request.Stderr, err)
		return ExitFailure
	}
	return writeInstallResult(parsed, request, items)
}

func installItems(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) ([]installItem, error) {
	if parsed.Command == cli.CommandInstallSource {
		return sourceInstallItems(parsed.Targets), nil
	}
	doc, err := readInstallLockfile(resolved.AgentsHome)
	if err != nil {
		return nil, err
	}
	if len(parsed.Targets) > 0 {
		return namedInstallItems(doc, parsed.Targets)
	}
	return changedInstallItems(ctx, doc, resolved, request)
}

func sourceInstallItems(targets []string) []installItem {
	items := make([]installItem, 0, len(targets))
	for _, target := range targets {
		items = append(items, installItem{source: target})
	}
	return items
}

func namedInstallItems(doc lockfile.Document, targets []string) ([]installItem, error) {
	items := make([]installItem, 0, len(targets))
	for _, name := range targets {
		entry, ok := doc.Skill(name)
		if !ok {
			return nil, fmt.Errorf("skill not found in global lockfile: %s", name)
		}
		items = append(items, installItem{name: name, source: installSourceURL(entry)})
	}
	return items, nil
}

func changedInstallItems(ctx context.Context, doc lockfile.Document, resolved config.Resolved, request Request) ([]installItem, error) {
	result, err := status.Check(ctx, gitRunner(request), status.Input{
		GitPath:    "git",
		AgentsHome: resolved.AgentsHome,
		CacheDir:   resolved.CacheDir,
		StateDir:   resolved.StateDir,
	})
	if err != nil {
		return nil, err
	}
	items := make([]installItem, 0, len(result.Statuses))
	for _, skillStatus := range result.Statuses {
		if skillStatus.Status != output.StatusUpdate && skillStatus.Status != output.StatusMissing {
			continue
		}
		entry, ok := doc.Skill(skillStatus.Name)
		if !ok {
			return nil, fmt.Errorf("source not found in lockfile: %s", skillStatus.Name)
		}
		items = append(items, installItem{name: skillStatus.Name, source: installSourceURL(entry)})
	}
	return items, nil
}

func runInstallItems(ctx context.Context, resolved config.Resolved, request Request, items []installItem) error {
	return withLockfilePreservation(ctx, resolved, request, func() error {
		for _, item := range items {
			if err := runSkillsAdd(ctx, resolved, request, item); err != nil {
				return err
			}
		}
		return nil
	})
}

func runSkillsAdd(ctx context.Context, resolved config.Resolved, request Request, item installItem) error {
	command, err := runner.ResolveSkillsCommand(runner.ResolveInput{Override: resolved.SkillsCommand})
	if err != nil {
		return err
	}
	args := append([]string(nil), command.Args...)
	args = append(args, "add", item.source, "-g", "-y", "--agent", "universal")
	if item.name != "" {
		args = append(args, "--skill", item.name)
	}
	result, err := toolRunner(request).Run(ctx, compare.Command{Name: command.Program, Args: args})
	if err != nil {
		return fmt.Errorf("skills add %s: %w", item.source, err)
	}
	if len(result.Stdout) > 0 {
		writeText(request.Stderr, string(result.Stdout))
	}
	return nil
}

func writeInstallNoop(parsed cli.Parsed, request Request) int {
	if parsed.Output == cli.OutputJSON {
		return writeSummary(request, parsed, nil)
	}
	if parsed.Output == cli.OutputJSONL {
		event := output.Event{Type: output.EventTypeSummary, OK: true}
		if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
		return ExitOK
	}
	if len(parsed.Targets) > 0 {
		writeText(request.Stdout, "OK      requested skills are up to date\n")
	} else {
		writeText(request.Stdout, "OK      all skills are up to date\n")
	}
	return ExitOK
}

func writeInstallResult(parsed cli.Parsed, request Request, items []installItem) int {
	if parsed.Output == cli.OutputJSON {
		return writeSummary(request, parsed, installStatuses(items))
	}
	if parsed.Output == cli.OutputJSONL {
		for _, item := range installStatuses(items) {
			event := output.Event{Type: output.EventTypeStatus, Name: item.Name, Status: item.Status, SourceURL: item.SourceURL}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
		return ExitOK
	}
	for _, item := range items {
		if item.name != "" {
			writeText(request.Stdout, fmt.Sprintf("INSTALL %s from %s\n", item.name, item.source))
		} else {
			writeText(request.Stdout, fmt.Sprintf("INSTALL %s\n", item.source))
		}
	}
	return ExitOK
}

func installStatuses(items []installItem) []output.SkillStatus {
	statuses := make([]output.SkillStatus, 0, len(items))
	for _, item := range items {
		statuses = append(statuses, output.SkillStatus{Name: item.name, Status: output.StatusOK, SourceURL: item.source})
	}
	return statuses
}

func readInstallLockfile(agentsHome string) (lockfile.Document, error) {
	raw, err := os.ReadFile(filepath.Join(agentsHome, ".skill-lock.json"))
	if err != nil {
		return lockfile.Document{}, fmt.Errorf("read lockfile: %w", err)
	}
	doc, err := lockfile.Parse(raw)
	if err != nil {
		return lockfile.Document{}, err
	}
	return doc, nil
}

func installSourceURL(entry lockfile.SkillEntry) string {
	if entry.SourceURL != "" {
		return entry.SourceURL
	}
	if entry.Source != "" {
		return fmt.Sprintf("https://github.com/%s.git", entry.Source)
	}
	return ""
}

func toolRunner(request Request) compare.CommandRunner {
	if request.ToolRunner != nil {
		return request.ToolRunner
	}
	return compare.ExecRunner{}
}
