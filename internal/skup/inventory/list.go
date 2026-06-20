package inventory

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func ListInstalled(agentsHome string) ([]string, error) {
	skillsDir := filepath.Join(agentsHome, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("read skills directory: %w", err)
	}

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	return names, nil
}
