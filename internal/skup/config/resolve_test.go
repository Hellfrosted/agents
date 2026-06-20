package config

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestResolve_usesExplicitPaths_whenOptionsProvided(t *testing.T) {
	// Given
	input := ResolveInput{
		Options: Options{
			AgentsHome: "/custom/agents",
			CacheDir:   "/custom/cache",
			StateDir:   "/custom/state",
		},
		Env: map[string]string{
			"SK_UP_AGENTS_HOME": "/env/agents",
			"SK_UP_CACHE_DIR":   "/env/cache",
			"SK_UP_STATE_DIR":   "/env/state",
		},
		Platform: Platform{GOOS: "linux", HomeDir: "/home/alice"},
	}

	// When
	got, err := Resolve(input)

	// Then
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.AgentsHome != "/custom/agents" {
		t.Fatalf("AgentsHome = %q, want %q", got.AgentsHome, "/custom/agents")
	}
	if got.CacheDir != "/custom/cache" {
		t.Fatalf("CacheDir = %q, want %q", got.CacheDir, "/custom/cache")
	}
	if got.StateDir != "/custom/state" {
		t.Fatalf("StateDir = %q, want %q", got.StateDir, "/custom/state")
	}
}

func TestResolve_usesEnvironment_whenOptionsMissing(t *testing.T) {
	// Given
	input := ResolveInput{
		Env: map[string]string{
			"SK_UP_AGENTS_HOME": "/env/agents",
			"SK_UP_CACHE_DIR":   "/env/cache",
			"SK_UP_STATE_DIR":   "/env/state",
		},
		Platform: Platform{GOOS: "linux", HomeDir: "/home/alice"},
	}

	// When
	got, err := Resolve(input)

	// Then
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.AgentsHome != "/env/agents" {
		t.Fatalf("AgentsHome = %q, want %q", got.AgentsHome, "/env/agents")
	}
	if got.CacheDir != "/env/cache" {
		t.Fatalf("CacheDir = %q, want %q", got.CacheDir, "/env/cache")
	}
	if got.StateDir != "/env/state" {
		t.Fatalf("StateDir = %q, want %q", got.StateDir, "/env/state")
	}
}

func TestResolve_defaultsDiffToolToZed_whenUnset(t *testing.T) {
	// Given
	input := ResolveInput{
		Options:  Options{AgentsHome: "/agents", CacheDir: "/cache", StateDir: "/state"},
		Platform: Platform{GOOS: "linux", HomeDir: "/home/alice"},
	}

	// When
	got, err := Resolve(input)

	// Then
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.DiffTool != "zed" {
		t.Fatalf("DiffTool = %q, want zed", got.DiffTool)
	}
}

func TestResolve_usesLinuxDefaults_whenXDGEnvironmentPresent(t *testing.T) {
	// Given
	input := ResolveInput{
		Env: map[string]string{
			"XDG_CACHE_HOME": "/xdg/cache",
			"XDG_STATE_HOME": "/xdg/state",
		},
		Platform: Platform{GOOS: "linux", HomeDir: "/home/alice"},
	}

	// When
	got, err := Resolve(input)

	// Then
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	wantAgents := filepath.Join("/home/alice", ".agents")
	if got.AgentsHome != wantAgents {
		t.Fatalf("AgentsHome = %q, want %q", got.AgentsHome, wantAgents)
	}
	if got.CacheDir != filepath.Join("/xdg/cache", "sk-up") {
		t.Fatalf("CacheDir = %q", got.CacheDir)
	}
	if got.StateDir != filepath.Join("/xdg/state", "sk-up") {
		t.Fatalf("StateDir = %q", got.StateDir)
	}
}

func TestResolve_usesWindowsDefaults_whenLocalAppDataPresent(t *testing.T) {
	// Given
	input := ResolveInput{
		Env: map[string]string{
			"LOCALAPPDATA": `C:\Users\alice\AppData\Local`,
			"USERPROFILE":  `C:\Users\alice`,
		},
		Platform: Platform{GOOS: "windows"},
	}

	// When
	got, err := Resolve(input)

	// Then
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got.AgentsHome != filepath.Join(`C:\Users\alice`, ".agents") {
		t.Fatalf("AgentsHome = %q", got.AgentsHome)
	}
	if got.CacheDir != filepath.Join(`C:\Users\alice\AppData\Local`, "sk-up", "cache") {
		t.Fatalf("CacheDir = %q", got.CacheDir)
	}
	if got.StateDir != filepath.Join(`C:\Users\alice\AppData\Local`, "sk-up", "state") {
		t.Fatalf("StateDir = %q", got.StateDir)
	}
}

func TestResolve_returnsMissingHomeError_whenNoAgentsHomeSourceExists(t *testing.T) {
	// Given
	input := ResolveInput{
		Platform: Platform{GOOS: "linux"},
	}

	// When
	_, err := Resolve(input)

	// Then
	if !errors.Is(err, ErrMissingHome) {
		t.Fatalf("Resolve error = %v, want ErrMissingHome", err)
	}
}

func TestResolve_returnsMissingHomeError_whenDefaultCacheNeedsHome(t *testing.T) {
	// Given
	input := ResolveInput{
		Options:  Options{AgentsHome: "/custom/agents"},
		Platform: Platform{GOOS: "linux"},
	}

	// When
	_, err := Resolve(input)

	// Then
	if !errors.Is(err, ErrMissingHome) {
		t.Fatalf("Resolve error = %v, want ErrMissingHome", err)
	}
}
