package status

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestCheck_usesRootTreeHash_whenSkillPathAtRepositoryRoot(t *testing.T) {
	// Given
	root := t.TempDir()
	agentsHome := filepath.Join(root, "agents")
	cacheDir := filepath.Join(root, "cache")
	writeStatusFile(t, filepath.Join(agentsHome, ".skill-lock.json"), `{
  "version": 1,
  "skills": {
    "goalcraft": {"sourceUrl": "https://github.com/example/root-skill.git", "skillPath": "SKILL.md"}
  }
}`)
	writeStatusFile(t, filepath.Join(agentsHome, "skills", "goalcraft", "SKILL.md"), "same\n")
	runner := &fakeGitRunner{
		archives: map[string][]byte{
			"HEAD": tarArchive(t, statusTarFile{name: "SKILL.md", contents: "same\n"}),
		},
		hashes: map[string]string{
			".": "root-tree-hash",
		},
	}

	// When
	got, err := Check(context.Background(), runner, Input{
		GitPath:    "git",
		AgentsHome: agentsHome,
		CacheDir:   cacheDir,
		StateDir:   filepath.Join(root, "state"),
	})

	// Then
	if err != nil {
		t.Fatalf("Check returned error: %v", err)
	}
	if len(got.Statuses) != 1 {
		t.Fatalf("Statuses = %#v", got.Statuses)
	}
	if got.Statuses[0].RemoteHash != "root-tree-hash" {
		t.Fatalf("RemoteHash = %q, want root-tree-hash", got.Statuses[0].RemoteHash)
	}
	if containsRevParseArg(runner.commands, "HEAD:.") {
		t.Fatal("Check used HEAD:. for a root-level skill path")
	}
}

func TestEnsureRepo_reclonesCache_whenExistingResetFails(t *testing.T) {
	// Given
	root := t.TempDir()
	repoDir := filepath.Join(root, "cache", "repo")
	writeStatusFile(t, filepath.Join(repoDir, ".git", "HEAD"), "ref: refs/heads/main\n")
	runner := &resetFailingRunner{}
	source := skillSource{
		name:      "react-doctor",
		sourceURL: "https://github.com/example/react-doctor.git",
		remoteDir: "skills/react-doctor",
		repoDir:   repoDir,
	}

	// When
	err := ensureRepo(context.Background(), runner, "git", source, nil)

	// Then
	if err != nil {
		t.Fatalf("ensureRepo returned error: %v", err)
	}
	if !containsGitSubcommand(runner.commands, "clone") {
		t.Fatalf("commands did not include fallback clone: %#v", runner.commands)
	}
	if !containsCommandArgSequenceAfter(runner.commands, "clone", "-c", "core.symlinks=false") {
		t.Fatalf("clone did not persist core.symlinks=false: %#v", runner.commands)
	}
	if !containsGitSubcommand(runner.commands, "sparse-checkout") {
		t.Fatalf("commands did not include sparse-checkout after fallback: %#v", runner.commands)
	}
	if !containsCommandArgSequence(runner.commands, "sparse-checkout", "-c", "core.symlinks=false") {
		t.Fatalf("sparse-checkout did not run with core.symlinks=false: %#v", runner.commands)
	}
}

func TestEnsureRepo_restoresStaleCache_whenFallbackSparseCheckoutFails(t *testing.T) {
	// Given
	root := t.TempDir()
	repoDir := filepath.Join(root, "cache", "repo")
	writeStatusFile(t, filepath.Join(repoDir, ".git", "HEAD"), "ref: refs/heads/main\n")
	writeStatusFile(t, filepath.Join(repoDir, "CACHE_MARKER"), "stale cache\n")
	runner := &resetThenSparseFailingRunner{}
	source := skillSource{
		name:      "react-doctor",
		sourceURL: "https://github.com/example/react-doctor.git",
		remoteDir: "skills/react-doctor",
		repoDir:   repoDir,
	}

	// When
	err := ensureRepo(context.Background(), runner, "git", source, nil)

	// Then
	if err == nil {
		t.Fatal("ensureRepo returned nil, want sparse-checkout error")
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, "CACHE_MARKER")); statErr != nil {
		t.Fatalf("stale cache marker missing after fallback failure: %v", statErr)
	}
}

func TestEnsureRepo_restoresStaleCache_whenFallbackCloneLeavesPartialRepo(t *testing.T) {
	// Given
	root := t.TempDir()
	repoDir := filepath.Join(root, "cache", "repo")
	writeStatusFile(t, filepath.Join(repoDir, ".git", "HEAD"), "ref: refs/heads/main\n")
	writeStatusFile(t, filepath.Join(repoDir, "CACHE_MARKER"), "stale cache\n")
	runner := &resetThenCloneFailingRunner{repoDir: repoDir}
	source := skillSource{
		name:      "react-doctor",
		sourceURL: "https://github.com/example/react-doctor.git",
		remoteDir: "skills/react-doctor",
		repoDir:   repoDir,
	}

	// When
	err := ensureRepo(context.Background(), runner, "git", source, nil)

	// Then
	if err == nil {
		t.Fatal("ensureRepo returned nil, want clone error")
	}
	if _, statErr := os.Stat(filepath.Join(repoDir, "CACHE_MARKER")); statErr != nil {
		t.Fatalf("stale cache marker missing after clone failure: %v", statErr)
	}
}
