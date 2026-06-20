package status

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/compare"
)

func containsRevParseArg(commands []compare.Command, want string) bool {
	for _, command := range commands {
		if containsArg(command.Args, "rev-parse") && command.Args[len(command.Args)-1] == want {
			return true
		}
	}
	return false
}

type resetFailingRunner struct {
	commands []compare.Command
}

func (r *resetFailingRunner) Run(_ context.Context, command compare.Command) (compare.CommandResult, error) {
	r.commands = append(r.commands, command)
	if containsGitSubcommand([]compare.Command{command}, "reset") {
		return compare.CommandResult{}, fmt.Errorf("reset failed")
	}
	return compare.CommandResult{}, nil
}

func containsGitSubcommand(commands []compare.Command, want string) bool {
	for _, command := range commands {
		if containsArg(command.Args, want) {
			return true
		}
	}
	return false
}

type resetThenSparseFailingRunner struct {
	commands []compare.Command
}

func (r *resetThenSparseFailingRunner) Run(_ context.Context, command compare.Command) (compare.CommandResult, error) {
	r.commands = append(r.commands, command)
	if containsArg(command.Args, "reset") {
		return compare.CommandResult{}, fmt.Errorf("reset failed")
	}
	if containsArg(command.Args, "sparse-checkout") {
		return compare.CommandResult{}, fmt.Errorf("sparse-checkout failed")
	}
	return compare.CommandResult{}, nil
}

type resetThenCloneFailingRunner struct {
	commands []compare.Command
	repoDir  string
}

func (r *resetThenCloneFailingRunner) Run(_ context.Context, command compare.Command) (compare.CommandResult, error) {
	r.commands = append(r.commands, command)
	if containsArg(command.Args, "reset") {
		return compare.CommandResult{}, fmt.Errorf("reset failed")
	}
	if containsArg(command.Args, "clone") {
		if err := os.MkdirAll(r.repoDir, 0o700); err != nil {
			return compare.CommandResult{}, err
		}
		if err := os.WriteFile(filepath.Join(r.repoDir, "PARTIAL_CLONE"), []byte("partial\n"), 0o600); err != nil {
			return compare.CommandResult{}, err
		}
		return compare.CommandResult{}, fmt.Errorf("clone failed")
	}
	return compare.CommandResult{}, nil
}

func containsCommandArgSequenceAfter(commands []compare.Command, anchor string, want ...string) bool {
	for _, command := range commands {
		anchorIndex := -1
		for i, arg := range command.Args {
			if arg == anchor {
				anchorIndex = i
				break
			}
		}
		if anchorIndex == -1 {
			continue
		}
		for i := anchorIndex + 1; i <= len(command.Args)-len(want); i++ {
			matched := true
			for j, arg := range want {
				if command.Args[i+j] != arg {
					matched = false
					break
				}
			}
			if matched {
				return true
			}
		}
	}
	return false
}

func containsCommandArgSequence(commands []compare.Command, anchor string, want ...string) bool {
	for _, command := range commands {
		if !containsArg(command.Args, anchor) {
			continue
		}
		for i := 0; i <= len(command.Args)-len(want); i++ {
			matched := true
			for j, arg := range want {
				if command.Args[i+j] != arg {
					matched = false
					break
				}
			}
			if matched {
				return true
			}
		}
	}
	return false
}
