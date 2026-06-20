package compare

import (
	"archive/tar"
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestCompareGitSkill_returnsOK_whenExportedTreeMatchesInstalledIgnoringCRLF(t *testing.T) {
	// Given
	installed := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "demo\r\n"})
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "skills/demo/SKILL.md", contents: "demo\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	runner := &fakeRunner{result: CommandResult{Stdout: archive.Bytes()}}
	input := SkillInput{
		GitPath:      "git",
		Repo:         "/repo",
		RemoteDir:    "skills/demo",
		InstalledDir: installed,
		ExportDir:    t.TempDir(),
	}

	// When
	got, err := CompareGitSkill(context.Background(), runner, input)

	// Then
	if err != nil {
		t.Fatalf("CompareGitSkill returned error: %v", err)
	}
	if got.Status != output.StatusOK {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusOK)
	}
	if len(runner.commands) != 1 {
		t.Fatalf("git command count = %d, want 1", len(runner.commands))
	}
}

func TestCompareGitSkill_returnsMissingWithoutGit_whenInstalledDirMissing(t *testing.T) {
	// Given
	root := t.TempDir()
	runner := &fakeRunner{}
	input := SkillInput{
		GitPath:      "git",
		Repo:         "/repo",
		RemoteDir:    "skills/demo",
		InstalledDir: filepath.Join(root, "missing"),
		ExportDir:    t.TempDir(),
	}

	// When
	got, err := CompareGitSkill(context.Background(), runner, input)

	// Then
	if err != nil {
		t.Fatalf("CompareGitSkill returned error: %v", err)
	}
	if got.Status != output.StatusMissing {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusMissing)
	}
	if len(runner.commands) != 0 {
		t.Fatalf("git command count = %d, want 0", len(runner.commands))
	}
}

func TestCompareGitSkill_returnsUpdate_whenExportedTreeDiffers(t *testing.T) {
	// Given
	installed := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "old\n"})
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "skills/demo/SKILL.md", contents: "new\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	runner := &fakeRunner{result: CommandResult{Stdout: archive.Bytes()}}
	input := SkillInput{GitPath: "git", Repo: "/repo", RemoteDir: "skills/demo", InstalledDir: installed, ExportDir: t.TempDir()}

	// When
	got, err := CompareGitSkill(context.Background(), runner, input)

	// Then
	if err != nil {
		t.Fatalf("CompareGitSkill returned error: %v", err)
	}
	if got.Status != output.StatusUpdate {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusUpdate)
	}
}

func TestCompareGitSkill_removesExistingExportDirBeforeExtraction(t *testing.T) {
	// Given
	installed := t.TempDir()
	exportDir := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "new\n"})
	writeTestFile(t, testFile{root: exportDir, name: "stale.md", contents: "stale\n"})
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "SKILL.md", contents: "new\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	// When
	_, err := CompareGitSkill(context.Background(), &fakeRunner{result: CommandResult{Stdout: archive.Bytes()}}, SkillInput{
		GitPath:      "git",
		Repo:         "/repo",
		RemoteDir:    ".",
		InstalledDir: installed,
		ExportDir:    exportDir,
	})

	// Then
	if err != nil {
		t.Fatalf("CompareGitSkill returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(exportDir, "stale.md")); !os.IsNotExist(err) {
		t.Fatalf("stale export file still exists or stat failed differently: %v", err)
	}
}
