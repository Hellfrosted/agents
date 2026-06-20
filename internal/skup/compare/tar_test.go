package compare

import (
	"archive/tar"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTar_writesFilesFromArchive(t *testing.T) {
	// Given
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "skills/demo/SKILL.md", contents: "demo\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}
	dest := t.TempDir()

	// When
	err := ExtractTar(&archive, dest)

	// Then
	if err != nil {
		t.Fatalf("ExtractTar returned error: %v", err)
	}
	raw, err := os.ReadFile(filepath.Join(dest, "skills", "demo", "SKILL.md"))
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(raw) != "demo\n" {
		t.Fatalf("file contents = %q", raw)
	}
}

func TestExtractTar_rejectsPathTraversal(t *testing.T) {
	// Given
	var archive bytes.Buffer
	tw := tar.NewWriter(&archive)
	writeTarFile(t, tarFile{name: "../escape.txt", contents: "bad\n"}, tw)
	if err := tw.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	// When
	err := ExtractTar(&archive, t.TempDir())

	// Then
	if !errors.Is(err, ErrUnsafeArchivePath) {
		t.Fatalf("ExtractTar error = %v, want ErrUnsafeArchivePath", err)
	}
}

type tarFile struct {
	name     string
	contents string
}

func writeTarFile(t *testing.T, file tarFile, tw *tar.Writer) {
	t.Helper()
	header := &tar.Header{
		Name: file.name,
		Mode: 0o600,
		Size: int64(len(file.contents)),
	}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatalf("WriteHeader returned error: %v", err)
	}
	if _, err := tw.Write([]byte(file.contents)); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
}
