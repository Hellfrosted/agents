package cli

import "errors"

var ErrUsage = errors.New("cli: usage error")

type Entrypoint string

const (
	EntrypointShort Entrypoint = "sk-up"
	EntrypointLong  Entrypoint = "skills-updates"
)

type Command string

const (
	CommandHelp          Command = "help"
	CommandList          Command = "list"
	CommandStatus        Command = "status"
	CommandDiff          Command = "diff"
	CommandOpenDiff      Command = "open-diff"
	CommandInstall       Command = "install"
	CommandInstallSource Command = "install-source"
	CommandSkip          Command = "skip"
	CommandUnskip        Command = "unskip"
	CommandSkips         Command = "skips"
	CommandRemove        Command = "remove"
)

type OutputMode string

const (
	OutputHuman OutputMode = "human"
	OutputJSON  OutputMode = "json"
	OutputJSONL OutputMode = "jsonl"
)

type Input struct {
	Argv0 string
	Args  []string
}

type Parsed struct {
	Entrypoint    Entrypoint
	Command       Command
	Targets       []string
	Output        OutputMode
	DryRun        bool
	AgentsHome    string
	CacheDir      string
	StateDir      string
	SkillsCommand string
	DiffTool      string
	Color         string
}
