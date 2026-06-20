package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/cli"
	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/config"
	"github.com/Hellfrosted/agents/internal/skup/lockfile"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/runner"
)

func executeRemove(ctx context.Context, parsed cli.Parsed, resolved config.Resolved, request Request) int {
	if err := withRemoveLock(ctx, resolved, request, func() error {
		for _, name := range parsed.Targets {
			if err := removeSkill(ctx, resolved, request, name); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		writeError(request.Stderr, err)
		return ExitFailure
	}
	return writeRemoveResult(parsed, request)
}

func withRemoveLock(ctx context.Context, resolved config.Resolved, request Request, action func() error) error {
	lockPath := filepath.Join(resolved.AgentsHome, ".skill-lock.json.lock")
	lock, err := lockfile.AcquireLock(ctx, lockPath, lockMetadata(request))
	if err != nil {
		if errors.Is(err, lockfile.ErrLockHeld) {
			return fmt.Errorf("%w: %s", lockfile.ErrLockTimeout, lockPath)
		}
		return err
	}
	defer func() {
		if err := lock.Release(); err != nil {
			writeError(request.Stderr, err)
		}
	}()
	return action()
}

func removeSkill(ctx context.Context, resolved config.Resolved, request Request, name string) error {
	if !safeSkillName(name) {
		return fmt.Errorf("invalid skill name for remove: %s", name)
	}
	if err := runSkillsRemove(ctx, resolved, request, name); err != nil {
		return err
	}
	if err := os.RemoveAll(filepath.Join(resolved.AgentsHome, "skills", name)); err != nil {
		return fmt.Errorf("remove installed skill directory: %w", err)
	}
	if err := removeLockEntry(resolved.AgentsHome, name); err != nil {
		return err
	}
	return removeSkipEntry(resolved.StateDir, name)
}

func runSkillsRemove(ctx context.Context, resolved config.Resolved, request Request, name string) error {
	command, err := runner.ResolveSkillsCommand(runner.ResolveInput{Override: resolved.SkillsCommand})
	if err != nil {
		return err
	}
	args := append([]string(nil), command.Args...)
	args = append(args, "remove", "-g", "-y", "--agent", "universal", "--skill", name)
	result, err := toolRunner(request).Run(ctx, compare.Command{Name: command.Program, Args: args})
	if err != nil {
		return fmt.Errorf("skills remove %s: %w", name, err)
	}
	if len(result.Stdout) > 0 {
		writeText(request.Stderr, string(result.Stdout))
	}
	return nil
}

func removeLockEntry(agentsHome string, name string) error {
	path := filepath.Join(agentsHome, ".skill-lock.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read lockfile: %w", err)
	}
	doc, err := lockfile.Parse(raw)
	if err != nil {
		return err
	}
	removed, err := doc.RemoveSkill(name)
	if err != nil {
		return err
	}
	if !removed {
		return nil
	}
	next, err := doc.Marshal()
	if err != nil {
		return err
	}
	return writeRawFile(path, next)
}

func removeSkipEntry(stateDir string, name string) error {
	skips, err := loadSkips(stateDir)
	if err != nil {
		return err
	}
	removed, err := skips.RemoveEntry(name)
	if err != nil {
		return err
	}
	if !removed {
		return nil
	}
	return saveSkips(stateDir, skips)
}

func writeRemoveResult(parsed cli.Parsed, request Request) int {
	statuses := make([]output.SkillStatus, 0, len(parsed.Targets))
	for _, name := range parsed.Targets {
		statuses = append(statuses, output.SkillStatus{Name: name, Status: output.StatusOK})
	}
	switch parsed.Output {
	case cli.OutputJSON:
		return writeSummary(request, parsed, statuses)
	case cli.OutputJSONL:
		for _, status := range statuses {
			event := output.Event{Type: output.EventTypeStatus, Name: status.Name, Status: status.Status}
			if err := output.WriteJSONLEvent(request.Stdout, event); err != nil {
				writeError(request.Stderr, err)
				return ExitUsage
			}
		}
	default:
		for _, name := range parsed.Targets {
			writeText(request.Stdout, fmt.Sprintf("REMOVE  %s\n", name))
		}
	}
	return ExitOK
}

func safeSkillName(name string) bool {
	return name != "" && name != "." && name != ".." && filepath.Base(name) == name
}
