package runner

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"unicode"
)

var (
	ErrInvalidCommand = errors.New("runner: invalid skills command")
	ErrNoSkillsRunner = errors.New("runner: no supported skills command found")
)

type LookupFunc func(name string) (string, bool)

type Command struct {
	Program string
	Args    []string
}

type ResolveInput struct {
	Override string
	Lookup   LookupFunc
}

func ResolveSkillsCommand(input ResolveInput) (Command, error) {
	if input.Override != "" {
		command, err := ParseCommand(input.Override)
		if err != nil {
			return Command{}, err
		}
		return command, nil
	}

	lookup := input.Lookup
	if lookup == nil {
		lookup = defaultLookup
	}

	for _, fallback := range fallbackCommands() {
		program, ok := lookup(fallback.Program)
		if !ok {
			continue
		}
		return Command{Program: program, Args: append([]string(nil), fallback.Args...)}, nil
	}

	return Command{}, ErrNoSkillsRunner
}

func ParseCommand(raw string) (Command, error) {
	tokens, err := splitCommand(raw)
	if err != nil {
		return Command{}, err
	}
	if len(tokens) == 0 {
		return Command{}, ErrInvalidCommand
	}
	return Command{Program: tokens[0], Args: append([]string(nil), tokens[1:]...)}, nil
}

func splitCommand(raw string) ([]string, error) {
	var tokens []string
	var current strings.Builder
	var quote rune

	for _, char := range raw {
		if quote != 0 {
			if char == quote {
				quote = 0
				continue
			}
			current.WriteRune(char)
			continue
		}

		switch {
		case char == '\'' || char == '"':
			quote = char
		case unicode.IsSpace(char):
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	if quote != 0 {
		return nil, fmt.Errorf("%w: unclosed quote", ErrInvalidCommand)
	}
	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}
	return tokens, nil
}

func fallbackCommands() []Command {
	return []Command{
		{Program: "pnpm", Args: []string{"dlx", "skills@latest"}},
		{Program: "bunx", Args: []string{"skills@latest"}},
		{Program: "deno", Args: []string{"run", "-A", "npm:skills@latest"}},
		{Program: "npx", Args: []string{"-y", "skills@latest"}},
	}
}

func defaultLookup(name string) (string, bool) {
	path, err := exec.LookPath(name)
	return path, err == nil
}
