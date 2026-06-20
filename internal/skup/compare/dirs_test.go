package compare

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Hellfrosted/agents/internal/skup/output"
)

func TestCompareDirs_returnsOK_whenFilesDifferOnlyByCRLF(t *testing.T) {
	// Given
	installed := t.TempDir()
	upstream := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "line one\r\nline two\r\n"})
	writeTestFile(t, testFile{root: upstream, name: "SKILL.md", contents: "line one\nline two\n"})

	// When
	got, err := CompareDirs(installed, upstream)

	// Then
	if err != nil {
		t.Fatalf("CompareDirs returned error: %v", err)
	}
	if got.Status != output.StatusOK {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusOK)
	}
}

func TestCompareDirs_returnsUpdate_whenFileContentDiffers(t *testing.T) {
	// Given
	installed := t.TempDir()
	upstream := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "old\n"})
	writeTestFile(t, testFile{root: upstream, name: "SKILL.md", contents: "new\n"})

	// When
	got, err := CompareDirs(installed, upstream)

	// Then
	if err != nil {
		t.Fatalf("CompareDirs returned error: %v", err)
	}
	if got.Status != output.StatusUpdate {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusUpdate)
	}
}

func TestCompareDirs_returnsMissing_whenInstalledDirectoryDoesNotExist(t *testing.T) {
	// Given
	root := t.TempDir()
	installed := filepath.Join(root, "missing")
	upstream := t.TempDir()

	// When
	got, err := CompareDirs(installed, upstream)

	// Then
	if err != nil {
		t.Fatalf("CompareDirs returned error: %v", err)
	}
	if got.Status != output.StatusMissing {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusMissing)
	}
}

func TestCompareDirs_returnsUpdate_whenFileSetDiffers(t *testing.T) {
	// Given
	installed := t.TempDir()
	upstream := t.TempDir()
	writeTestFile(t, testFile{root: installed, name: "SKILL.md", contents: "same\n"})
	writeTestFile(t, testFile{root: installed, name: "local-only.md", contents: "extra\n"})
	writeTestFile(t, testFile{root: upstream, name: "SKILL.md", contents: "same\n"})

	// When
	got, err := CompareDirs(installed, upstream)

	// Then
	if err != nil {
		t.Fatalf("CompareDirs returned error: %v", err)
	}
	if got.Status != output.StatusUpdate {
		t.Fatalf("Status = %q, want %q", got.Status, output.StatusUpdate)
	}
}

type testFile struct {
	root     string
	name     string
	contents string
}

func writeTestFile(t *testing.T, file testFile) {
	t.Helper()
	path := filepath.Join(file.root, file.name)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}
	if err := os.WriteFile(path, []byte(file.contents), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}
