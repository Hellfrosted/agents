package status

import (
	"crypto/sha256"
	"fmt"
	"path"
	"path/filepath"

	"github.com/Hellfrosted/agents/internal/skup/lockfile"
)

type skillSource struct {
	name      string
	sourceURL string
	remoteDir string
	repoDir   string
	exportDir string
	installed string
}

func (s skillSource) compareDir() string {
	if s.remoteDir == "" || s.remoteDir == "." {
		return s.exportDir
	}
	return filepath.Join(s.exportDir, filepath.FromSlash(s.remoteDir))
}

func newSkillSource(name string, entry lockfile.SkillEntry, input Input) skillSource {
	sourceURL := sourceURLFor(entry)
	return skillSource{
		name:      name,
		sourceURL: sourceURL,
		remoteDir: remoteDirFor(entry.SkillPath),
		repoDir:   filepath.Join(input.CacheDir, "repos", cacheKey(sourceURL)),
		exportDir: filepath.Join(input.CacheDir, "exports", name),
		installed: filepath.Join(input.AgentsHome, "skills", name),
	}
}

func sourceURLFor(entry lockfile.SkillEntry) string {
	if entry.SourceURL != "" {
		return entry.SourceURL
	}
	if entry.Source != "" {
		return fmt.Sprintf("https://github.com/%s.git", entry.Source)
	}
	return ""
}

func remoteDirFor(skillPath string) string {
	dir := path.Dir(path.Clean(skillPath))
	if dir == "." || dir == "/" {
		return "."
	}
	return dir
}

func cacheKey(value string) string {
	sum := sha256.Sum256([]byte(value))
	return fmt.Sprintf("%x", sum)
}
