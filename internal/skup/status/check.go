package status

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/compare"
	"github.com/Hellfrosted/agents/internal/skup/lockfile"
	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/state"
)

func Check(ctx context.Context, runner compare.CommandRunner, input Input) (Result, error) {
	doc, err := readLockfile(input.AgentsHome)
	if err != nil {
		return Result{}, err
	}
	skips, err := readSkips(input.StateDir)
	if err != nil {
		return Result{}, err
	}

	statuses := make([]output.SkillStatus, 0, len(doc.SkillNames()))
	for _, name := range doc.SkillNames() {
		if !targetSelected(name, input.Targets) {
			continue
		}
		status, err := checkSkill(ctx, runner, input, doc, skips, name)
		if err != nil {
			return Result{}, err
		}
		statuses = append(statuses, status)
	}
	return Result{Statuses: statuses}, nil
}

func checkSkill(ctx context.Context, runner compare.CommandRunner, input Input, doc lockfile.Document, skips state.Skips, name string) (output.SkillStatus, error) {
	entry, ok := doc.Skill(name)
	if !ok {
		return output.SkillStatus{}, fmt.Errorf("skill %s not found in lockfile", name)
	}
	source := newSkillSource(name, entry, input)
	if err := ensureRepo(ctx, runner, input.GitPath, source); err != nil {
		return output.SkillStatus{}, err
	}
	hash, err := remoteHash(ctx, runner, input.GitPath, source)
	if err != nil {
		return output.SkillStatus{}, err
	}
	result, err := compare.CompareGitSkill(ctx, runner, compare.SkillInput{
		GitPath:      input.GitPath,
		Repo:         source.repoDir,
		RemoteDir:    source.remoteDir,
		InstalledDir: source.installed,
		ExportDir:    source.exportDir,
	})
	if err != nil {
		return output.SkillStatus{}, err
	}
	return statusFor(name, source, hash, result.Status, skips), nil
}

func readLockfile(agentsHome string) (lockfile.Document, error) {
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

func readSkips(stateDir string) (state.Skips, error) {
	raw, err := os.ReadFile(filepath.Join(stateDir, "skips.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return state.NewSkips(), nil
		}
		return state.Skips{}, fmt.Errorf("read skips: %w", err)
	}
	return state.ParseSkips(raw)
}

func statusFor(name string, source skillSource, hash string, status output.Status, skips state.Skips) output.SkillStatus {
	if status == output.StatusUpdate {
		if skip, ok := skips.Entry(name); ok && skip.RemoteHash == hash {
			status = output.StatusSkipped
		}
	}
	return output.SkillStatus{
		Name:         name,
		Status:       status,
		SourceURL:    source.sourceURL,
		RemoteHash:   hash,
		InstalledDir: source.installed,
		CompareDir:   source.compareDir(),
	}
}

func targetSelected(name string, targets []string) bool {
	if len(targets) == 0 {
		return true
	}
	for _, target := range targets {
		if target == name {
			return true
		}
	}
	return false
}
