package status

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestCheck_filtersLargeTargetList_whenTargetsProvided(t *testing.T) {
	// Given
	root := t.TempDir()
	fixture := newLargeTargetFilterFixture(t, root)
	runner := fixture.runner()

	// When
	got, err := Check(context.Background(), runner, fixture.input())

	// Then
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	requireLargeTargetStatuses(t, got.Statuses, fixture)
	if got := countCommandWithArg(runner.commands, "clone"); got != len(fixture.selected) {
		t.Fatalf("clone command count = %d, want %d", got, len(fixture.selected))
	}
	if got := countCommandWithArg(runner.commands, "rev-parse"); got != len(fixture.selected) {
		t.Fatalf("rev-parse command count = %d, want %d", got, len(fixture.selected))
	}
	if got := countCommandWithArg(runner.commands, "archive"); got != 0 {
		t.Fatalf("archive command count = %d, want 0", got)
	}
}

func BenchmarkCheck_filtersLargeTargetList(b *testing.B) {
	root := b.TempDir()
	fixture := newLargeTargetFilterFixture(b, root)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner := fixture.runner()
		got, err := Check(context.Background(), runner, fixture.input())
		if err != nil {
			b.Fatalf("Check returned error: %v", err)
		}
		requireLargeTargetStatuses(b, got.Statuses, fixture)
	}
}

func BenchmarkCheck_filtersSingleTarget(b *testing.B) {
	root := b.TempDir()
	fixture := newLargeTargetFilterFixture(b, root)
	fixture.targets = []string{fixture.selected[0]}
	fixture.selected = fixture.selected[:1]

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		runner := fixture.runner()
		got, err := Check(context.Background(), runner, fixture.input())
		if err != nil {
			b.Fatalf("Check returned error: %v", err)
		}
		requireLargeTargetStatuses(b, got.Statuses, fixture)
	}
}

func BenchmarkTargetFilter_selectsLargeTargetList(b *testing.B) {
	targets := newTargetFilter(largeTargetList([]string{"skill-1023"}, 2048))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !targets.selected("skill-1023") {
			b.Fatal("target not selected")
		}
	}
}

func BenchmarkTargetFilter_selectsSingleTarget(b *testing.B) {
	targets := newTargetFilter([]string{"skill-1023"})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if !targets.selected("skill-1023") {
			b.Fatal("target not selected")
		}
	}
}

type largeTargetFilterFixture struct {
	agentsHome string
	cacheDir   string
	stateDir   string
	targets    []string
	selected   []string
	hashes     map[string]string
}

func newLargeTargetFilterFixture(tb testing.TB, root string) largeTargetFilterFixture {
	tb.Helper()
	const skillCount = 1024
	selected := []string{"skill-0007", "skill-0512", "skill-1023"}
	agentsHome := filepath.Join(root, "agents")
	writeStatusFileTB(tb, filepath.Join(agentsHome, ".skill-lock.json"), largeTargetLockfile(skillCount))
	return largeTargetFilterFixture{
		agentsHome: agentsHome,
		cacheDir:   filepath.Join(root, "cache"),
		stateDir:   filepath.Join(root, "state"),
		targets:    largeTargetList(selected, 2048),
		selected:   selected,
		hashes: map[string]string{
			"skills/skill-0007": "hash-skill-0007",
			"skills/skill-0512": "hash-skill-0512",
			"skills/skill-1023": "hash-skill-1023",
		},
	}
}

func writeStatusFileTB(tb testing.TB, path string, contents string) {
	tb.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		tb.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		tb.Fatalf("WriteFile returned error: %v", err)
	}
}

func (f largeTargetFilterFixture) input() Input {
	return Input{
		GitPath:    "git",
		AgentsHome: f.agentsHome,
		CacheDir:   f.cacheDir,
		StateDir:   f.stateDir,
		Targets:    f.targets,
	}
}

func (f largeTargetFilterFixture) runner() *fakeGitRunner {
	return &fakeGitRunner{hashes: f.hashes}
}

func largeTargetLockfile(skillCount int) string {
	var buffer bytes.Buffer
	buffer.WriteString(`{"version":1,"skills":{`)
	for i := 0; i < skillCount; i++ {
		if i > 0 {
			buffer.WriteByte(',')
		}
		name := fmt.Sprintf("skill-%04d", i)
		fmt.Fprintf(
			&buffer,
			"%q:{%q:%q,%q:%q}",
			name,
			"sourceUrl",
			fmt.Sprintf("https://github.com/example/%s.git", name),
			"skillPath",
			fmt.Sprintf("skills/%s/SKILL.md", name),
		)
	}
	buffer.WriteString(`}}`)
	return buffer.String()
}

func largeTargetList(selected []string, missCount int) []string {
	targets := make([]string, 0, missCount+len(selected))
	for i := 0; i < missCount; i++ {
		targets = append(targets, fmt.Sprintf("missing-%04d", i))
	}
	targets = append(targets, selected...)
	return targets
}

func requireLargeTargetStatuses(tb testing.TB, got []output.SkillStatus, fixture largeTargetFilterFixture) {
	tb.Helper()
	if len(got) != len(fixture.selected) {
		tb.Fatalf("Statuses len = %d, want %d; statuses=%#v", len(got), len(fixture.selected), got)
	}
	for i, name := range fixture.selected {
		status := got[i]
		if status.Name != name {
			tb.Fatalf("Statuses[%d].Name = %q, want %q; statuses=%#v", i, status.Name, name, got)
		}
		if status.Status != output.StatusMissing {
			tb.Fatalf("Statuses[%d].Status = %q, want %q", i, status.Status, output.StatusMissing)
		}
		if status.RemoteHash != fixture.hashes[fmt.Sprintf("skills/%s", name)] {
			tb.Fatalf("Statuses[%d].RemoteHash = %q", i, status.RemoteHash)
		}
	}
}
