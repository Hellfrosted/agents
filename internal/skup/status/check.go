package status

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

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

	sources := make([]skillSource, 0, len(doc.SkillNames()))
	for _, name := range doc.SkillNames() {
		if !targetSelected(name, input.Targets) {
			continue
		}
		source, err := skillSourceFor(input, doc, name)
		if err != nil {
			return Result{}, err
		}
		sources = append(sources, source)
	}

	if err := ensureSources(ctx, runner, input, sources); err != nil {
		return Result{}, err
	}
	statuses, err := checkSources(ctx, runner, input, skips, sources)
	if err != nil {
		return Result{}, err
	}
	return Result{Statuses: statuses}, nil
}

func skillSourceFor(input Input, doc lockfile.Document, name string) (skillSource, error) {
	entry, ok := doc.Skill(name)
	if !ok {
		return skillSource{}, fmt.Errorf("skill %s not found in lockfile", name)
	}
	return newSkillSource(name, entry, input), nil
}

func ensureSources(ctx context.Context, runner compare.CommandRunner, input Input, sources []skillSource) error {
	unique := uniqueRepoSources(sources)
	jobs := make(chan skillSource)
	errs := make(chan error, 1)
	var wg sync.WaitGroup

	for range statusWorkers(len(unique)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for source := range jobs {
				if err := ensureRepo(ctx, runner, input.GitPath, source); err != nil {
					select {
					case errs <- err:
					default:
					}
				}
			}
		}()
	}
	for _, source := range unique {
		select {
		case jobs <- source:
		case err := <-errs:
			close(jobs)
			wg.Wait()
			return err
		}
	}
	close(jobs)
	wg.Wait()
	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func uniqueRepoSources(sources []skillSource) []skillSource {
	seen := make(map[string]struct{}, len(sources))
	unique := make([]skillSource, 0, len(sources))
	for _, source := range sources {
		if _, ok := seen[source.repoDir]; ok {
			continue
		}
		seen[source.repoDir] = struct{}{}
		unique = append(unique, source)
	}
	return unique
}

func checkSources(ctx context.Context, runner compare.CommandRunner, input Input, skips state.Skips, sources []skillSource) ([]output.SkillStatus, error) {
	statuses := make([]output.SkillStatus, len(sources))
	jobs := make(chan int)
	errs := make(chan error, 1)
	var wg sync.WaitGroup

	for range statusWorkers(len(sources)) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				status, err := checkSource(ctx, runner, input, skips, sources[index])
				if err != nil {
					select {
					case errs <- err:
					default:
					}
					continue
				}
				statuses[index] = status
			}
		}()
	}
	for index := range sources {
		select {
		case jobs <- index:
		case err := <-errs:
			close(jobs)
			wg.Wait()
			return nil, err
		}
	}
	close(jobs)
	wg.Wait()
	select {
	case err := <-errs:
		return nil, err
	default:
		return statuses, nil
	}
}

func checkSource(ctx context.Context, runner compare.CommandRunner, input Input, skips state.Skips, source skillSource) (output.SkillStatus, error) {
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
	return statusFor(source.name, source, hash, result.Status, skips), nil
}

func statusWorkers(count int) int {
	if count < 2 {
		return count
	}
	workers := runtime.NumCPU()
	if workers < 2 {
		workers = 2
	}
	if workers > 8 {
		workers = 8
	}
	if count < workers {
		return count
	}
	return workers
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
