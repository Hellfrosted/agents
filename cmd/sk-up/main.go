package main

import (
	"context"
	"os"

	"github.com/Hellfrosted/agents/internal/skup/app"
)

func main() {
	env := envMap(os.Environ())
	code := app.Execute(context.Background(), app.Request{
		Argv0:  argv0ForEnv(os.Args[0], env),
		Args:   os.Args[1:],
		Env:    env,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	os.Exit(code)
}

func argv0ForEnv(argv0 string, env map[string]string) string {
	if env["SK_UP_ENTRYPOINT"] == "skills-updates" {
		return "skills-updates"
	}
	return argv0
}

func envMap(entries []string) map[string]string {
	env := make(map[string]string, len(entries))
	for _, entry := range entries {
		for i, char := range entry {
			if char == '=' {
				env[entry[:i]] = entry[i+1:]
				break
			}
		}
	}
	return env
}
