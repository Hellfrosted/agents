package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
)

var (
	ErrInvalidSkillName = errors.New("state: invalid skill name")
	ErrInvalidSkips     = errors.New("state: invalid skips document")
)

type SkipEntry struct {
	RemoteHash string `json:"remoteHash"`
	SourceURL  string `json:"sourceUrl"`
	SkippedAt  string `json:"skippedAt"`
}

type Skips struct {
	entries map[string]SkipEntry
}

func NewSkips() Skips {
	return Skips{entries: make(map[string]SkipEntry)}
}

func ParseSkips(raw []byte) (Skips, error) {
	if len(raw) == 0 {
		return NewSkips(), nil
	}

	var decoded struct {
		Skips map[string]SkipEntry `json:"skips"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return Skips{}, fmt.Errorf("%w: %w", ErrInvalidSkips, err)
	}
	if decoded.Skips == nil {
		return NewSkips(), nil
	}

	entries := make(map[string]SkipEntry, len(decoded.Skips))
	for name, entry := range decoded.Skips {
		entries[name] = entry
	}
	return Skips{entries: entries}, nil
}

func (s Skips) Entry(name string) (SkipEntry, bool) {
	entry, ok := s.entries[name]
	return entry, ok
}

func (s Skips) Names() []string {
	names := make([]string, 0, len(s.entries))
	for name := range s.entries {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (s Skips) SetEntry(name string, entry SkipEntry) error {
	if !validSkillName(name) {
		return ErrInvalidSkillName
	}
	s.entries[name] = entry
	return nil
}

func (s Skips) RemoveEntry(name string) (bool, error) {
	if !validSkillName(name) {
		return false, ErrInvalidSkillName
	}
	if _, ok := s.entries[name]; !ok {
		return false, nil
	}
	delete(s.entries, name)
	return true, nil
}

func (s Skips) Marshal() ([]byte, error) {
	raw, err := json.MarshalIndent(struct {
		Skips map[string]SkipEntry `json:"skips"`
	}{Skips: s.entries}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal skips: %w", err)
	}
	return append(raw, '\n'), nil
}

func validSkillName(name string) bool {
	return name != "" && name != "." && name != ".." && filepath.Base(name) == name
}
