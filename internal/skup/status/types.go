package status

import "github.com/Hellfrosted/agents/internal/skup/output"

type Input struct {
	GitPath    string
	AgentsHome string
	CacheDir   string
	StateDir   string
	Targets    []string
	Progress   ProgressFunc
}

type Result struct {
	Statuses []output.SkillStatus
}

type ProgressFunc func(output.Event)
