package status

import "github.com/Hellfrosted/agents/internal/skup/output"

type Input struct {
	GitPath    string
	AgentsHome string
	CacheDir   string
	StateDir   string
	Targets    []string
}

type Result struct {
	Statuses []output.SkillStatus
}
