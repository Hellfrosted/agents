package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/state"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

func executeSkip(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
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
	if len(result.Statuses) == 0 {
		writeError(request.Stderr, fmt.Errorf("skill %s not found in lockfile", parsed.Targets[0]))
		return ExitUsage
	}

	skillStatus := result.Statuses[0]
	if skillStatus.Status == output.StatusOK || skillStatus.Status == output.StatusSkipped {
		return writeSkipNoop(parsed, request, skillStatus)
	}
	if skillStatus.RemoteHash == "" {
		writeError(request.Stderr, fmt.Errorf("skill %s has no upstream hash to skip", skillStatus.Name))
		return ExitUsage
	}

	skips, err := loadSkips(resolved.StateDir)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	err = skips.SetEntry(skillStatus.Name, state.SkipEntry{
		RemoteHash: skillStatus.RemoteHash,
		SourceURL:  skillStatus.SourceURL,
		SkippedAt:  nowUTC(request).Format(time.RFC3339Nano),
	})
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	if err := saveSkips(resolved.StateDir, skips); err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	return writeSkipResult(parsed, request, skillStatus)
}

func executeUnskip(parsed cli.Parsed, resolved config.Resolved, request Request) int {
	skips, err := loadSkips(resolved.StateDir)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	removed, err := skips.RemoveEntry(parsed.Targets[0])
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	if removed {
		if err := saveSkips(resolved.StateDir, skips); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	}
	return writeUnskipResult(parsed, request, parsed.Targets[0], removed)
}

func executeSkips(parsed cli.Parsed, resolved config.Resolved, request Request) int {
	skips, err := loadSkips(resolved.StateDir)
	if err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	statuses := skipStatuses(skips)
	switch parsed.Output {
	case cli.OutputJSON:
		summary := output.Summary{
			OK:         true,
			Command:    string(parsed.Command),
			Entrypoint: string(parsed.Entrypoint),
			Statuses:   statuses,
		}
		if err := output.WriteSummaryJSON(request.Stdout, summary); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	case cli.OutputJSONL:
		for _, status := range statuses {
			event := output.Event{Type: output.EventTypeStatus, Name: status.Name, Status: status.Status, SourceURL: status.SourceURL, RemoteHash: status.RemoteHash}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
	default:
		for _, status := range statuses {
			writeText(request.Stdout, fmt.Sprintf("%-7s %s %s\n", strings.ToUpper(string(status.Status)), status.Name, status.RemoteHash))
		}
	}
	return ExitOK
}

func writeSkipNoop(parsed cli.Parsed, request Request, skillStatus output.SkillStatus) int {
	switch parsed.Output {
	case cli.OutputJSON:
		return writeSummary(request, parsed, []output.SkillStatus{skillStatus})
	case cli.OutputJSONL:
		event := output.Event{Type: output.EventTypeStatus, Name: skillStatus.Name, Status: skillStatus.Status, SourceURL: skillStatus.SourceURL, RemoteHash: skillStatus.RemoteHash}
		if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	default:
		writeText(request.Stdout, fmt.Sprintf("OK      %s has no current update to skip\n", skillStatus.Name))
	}
	return ExitOK
}

func writeSkipResult(parsed cli.Parsed, request Request, skillStatus output.SkillStatus) int {
	skillStatus.Status = output.StatusSkipped
	switch parsed.Output {
	case cli.OutputJSON:
		return writeSummary(request, parsed, []output.SkillStatus{skillStatus})
	case cli.OutputJSONL:
		event := output.Event{Type: output.EventTypeStatus, Name: skillStatus.Name, Status: output.StatusSkipped, SourceURL: skillStatus.SourceURL, RemoteHash: skillStatus.RemoteHash}
		if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	default:
		writeText(request.Stdout, fmt.Sprintf("SKIPPED %s\n", skillStatus.Name))
	}
	return ExitOK
}

func writeUnskipResult(parsed cli.Parsed, request Request, name string, removed bool) int {
	status := output.SkillStatus{Name: name, Status: output.StatusOK}
	switch parsed.Output {
	case cli.OutputJSON:
		return writeSummary(request, parsed, []output.SkillStatus{status})
	case cli.OutputJSONL:
		event := output.Event{Type: output.EventTypeStatus, Name: name, Status: output.StatusOK}
		if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
			writeError(request.Stderr, err)
			return ExitUsage
		}
	default:
		if removed {
			writeText(request.Stdout, fmt.Sprintf("OK      removed saved skip for %s\n", name))
		} else {
			writeText(request.Stdout, fmt.Sprintf("OK      no saved skip for %s\n", name))
		}
	}
	return ExitOK
}

func writeSummary(request Request, parsed cli.Parsed, statuses []output.SkillStatus) int {
	summary := output.Summary{
		OK:         true,
		Command:    string(parsed.Command),
		Entrypoint: string(parsed.Entrypoint),
		Statuses:   statuses,
	}
	if err := output.WriteSummaryJSON(request.Stdout, summary); err != nil {
		writeError(request.Stderr, err)
		return ExitUsage
	}
	return ExitOK
}

func skipStatuses(skips state.Skips) []output.SkillStatus {
	names := skips.Names()
	statuses := make([]output.SkillStatus, 0, len(names))
	for _, name := range names {
		entry, _ := skips.Entry(name)
		statuses = append(statuses, output.SkillStatus{
			Name:       name,
			Status:     output.StatusSkipped,
			SourceURL:  entry.SourceURL,
			RemoteHash: entry.RemoteHash,
		})
	}
	return statuses
}

func loadSkips(stateDir string) (state.Skips, error) {
	raw, err := os.ReadFile(filepath.Join(stateDir, "skips.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return state.NewSkips(), nil
		}
		return state.Skips{}, fmt.Errorf("read skips: %w", err)
	}
	skips, err := state.ParseSkips(raw)
	if err != nil {
		return state.Skips{}, err
	}
	return skips, nil
}

func saveSkips(stateDir string, skips state.Skips) error {
	raw, err := skips.Marshal()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return fmt.Errorf("create state dir: %w", err)
	}
	path := filepath.Join(stateDir, "skips.json")
	tmp, err := os.CreateTemp(stateDir, "skips-*.json")
	if err != nil {
		return fmt.Errorf("create skips temp file: %w", err)
	}
	tmpPath := tmp.Name()
	cleanup := true
	defer func() {
		if cleanup {
			_ = os.Remove(tmpPath)
		}
	}()
	if _, err := tmp.Write(raw); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write skips temp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close skips temp file: %w", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("replace skips file: %w", err)
	}
	cleanup = false
	return nil
}

func nowUTC(request Request) time.Time {
	if request.Now != nil {
		return request.Now().UTC()
	}
	return time.Now().UTC()
}
