package config

import (
	"errors"
	"path/filepath"
	"runtime"
)

const appName = "sk-up"

var ErrMissingHome = errors.New("config: missing home directory")

type Options struct {
	AgentsHome    string
	CacheDir      string
	StateDir      string
	SkillsCommand string
	DiffTool      string
	Color         string
}

type Platform struct {
	GOOS    string
	HomeDir string
}

type ResolveInput struct {
	Options  Options
	Env      map[string]string
	Platform Platform
}

type Resolved struct {
	AgentsHome    string
	CacheDir      string
	StateDir      string
	SkillsCommand string
	DiffTool      string
	Color         string
}

func Resolve(input ResolveInput) (Resolved, error) {
	platform := normalizePlatform(input.Platform, input.Env)

	agentsHome, err := resolveAgentsHome(input.Options, input.Env, platform)
	if err != nil {
		return Resolved{}, err
	}

	cacheDir, err := resolveCacheDir(input.Options, input.Env, platform)
	if err != nil {
		return Resolved{}, err
	}

	stateDir, err := resolveStateDir(input.Options, input.Env, platform)
	if err != nil {
		return Resolved{}, err
	}

	return Resolved{
		AgentsHome:    agentsHome,
		CacheDir:      cacheDir,
		StateDir:      stateDir,
		SkillsCommand: firstNonEmpty(input.Options.SkillsCommand, envValue(input.Env, "SK_UP_SKILLS_COMMAND")),
		DiffTool:      firstNonEmpty(input.Options.DiffTool, envValue(input.Env, "SK_UP_DIFF_TOOL"), "zed"),
		Color:         firstNonEmpty(input.Options.Color, envValue(input.Env, "SK_UP_COLOR"), "auto"),
	}, nil
}

func normalizePlatform(platform Platform, env map[string]string) Platform {
	goos := firstNonEmpty(platform.GOOS, runtime.GOOS)
	homeDir := firstNonEmpty(platform.HomeDir, envValue(env, "HOME"), envValue(env, "USERPROFILE"))
	return Platform{GOOS: goos, HomeDir: homeDir}
}

func resolveAgentsHome(options Options, env map[string]string, platform Platform) (string, error) {
	agentsHome := firstNonEmpty(
		options.AgentsHome,
		envValue(env, "SK_UP_AGENTS_HOME"),
		envValue(env, "AGENTS_HOME"),
	)
	if agentsHome != "" {
		return agentsHome, nil
	}
	if platform.HomeDir == "" {
		return "", ErrMissingHome
	}
	return filepath.Join(platform.HomeDir, ".agents"), nil
}

func resolveCacheDir(options Options, env map[string]string, platform Platform) (string, error) {
	cacheDir := firstNonEmpty(options.CacheDir, envValue(env, "SK_UP_CACHE_DIR"))
	if cacheDir != "" {
		return cacheDir, nil
	}

	switch platform.GOOS {
	case "windows":
		base := envValue(env, "LOCALAPPDATA")
		if base == "" && platform.HomeDir != "" {
			base = filepath.Join(platform.HomeDir, "AppData", "Local")
		}
		if base == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(base, appName, "cache"), nil
	case "darwin":
		if platform.HomeDir == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(platform.HomeDir, "Library", "Caches", appName), nil
	default:
		base := envValue(env, "XDG_CACHE_HOME")
		if base == "" && platform.HomeDir != "" {
			base = filepath.Join(platform.HomeDir, ".cache")
		}
		if base == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(base, appName), nil
	}
}

func resolveStateDir(options Options, env map[string]string, platform Platform) (string, error) {
	stateDir := firstNonEmpty(options.StateDir, envValue(env, "SK_UP_STATE_DIR"))
	if stateDir != "" {
		return stateDir, nil
	}

	switch platform.GOOS {
	case "windows":
		base := envValue(env, "LOCALAPPDATA")
		if base == "" && platform.HomeDir != "" {
			base = filepath.Join(platform.HomeDir, "AppData", "Local")
		}
		if base == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(base, appName, "state"), nil
	case "darwin":
		if platform.HomeDir == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(platform.HomeDir, "Library", "Application Support", appName), nil
	default:
		base := envValue(env, "XDG_STATE_HOME")
		if base == "" && platform.HomeDir != "" {
			base = filepath.Join(platform.HomeDir, ".local", "state")
		}
		if base == "" {
			return "", ErrMissingHome
		}
		return filepath.Join(base, appName), nil
	}
}

func envValue(env map[string]string, name string) string {
	if env == nil {
		return ""
	}
	return env[name]
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
