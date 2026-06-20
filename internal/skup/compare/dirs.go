package compare

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/output"
)

type Result struct {
	Status output.Status
}

func CompareDirs(installedDir string, upstreamDir string) (Result, error) {
	if _, err := os.Stat(installedDir); err != nil {
		if os.IsNotExist(err) {
			return Result{Status: output.StatusMissing}, nil
		}
		return Result{}, fmt.Errorf("stat installed directory: %w", err)
	}

	installedFiles, err := collectFiles(installedDir)
	if err != nil {
		return Result{}, fmt.Errorf("collect installed files: %w", err)
	}
	upstreamFiles, err := collectFiles(upstreamDir)
	if err != nil {
		return Result{}, fmt.Errorf("collect upstream files: %w", err)
	}

	if len(installedFiles) != len(upstreamFiles) {
		return Result{Status: output.StatusUpdate}, nil
	}

	for name, installed := range installedFiles {
		upstream, ok := upstreamFiles[name]
		if !ok {
			return Result{Status: output.StatusUpdate}, nil
		}
		if !bytes.Equal(normalizeLineEndings(installed), normalizeLineEndings(upstream)) {
			return Result{Status: output.StatusUpdate}, nil
		}
	}

	return Result{Status: output.StatusOK}, nil
}

func collectFiles(root string) (map[string][]byte, error) {
	files := make(map[string][]byte)
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return fmt.Errorf("relative path: %w", err)
		}
		raw, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		files[filepath.ToSlash(rel)] = raw
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func normalizeLineEndings(raw []byte) []byte {
	return bytes.ReplaceAll(raw, []byte("\r\n"), []byte("\n"))
}
