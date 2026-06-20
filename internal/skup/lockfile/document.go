package lockfile

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
)

var (
	ErrInvalidSkillName = errors.New("lockfile: invalid skill name")
	ErrInvalidDocument  = errors.New("lockfile: invalid document")
)

type SkillEntry struct {
	SourceURL  string
	Source     string
	SkillPath  string
	PluginName string
	Raw        json.RawMessage
}

type Document struct {
	fields map[string]json.RawMessage
	skills map[string]json.RawMessage
}

func Parse(raw []byte) (Document, error) {
	fields := make(map[string]json.RawMessage)
	if err := json.Unmarshal(raw, &fields); err != nil {
		return Document{}, fmt.Errorf("%w: %w", ErrInvalidDocument, err)
	}

	skillsRaw, ok := fields["skills"]
	if !ok {
		return Document{}, fmt.Errorf("%w: missing skills", ErrInvalidDocument)
	}

	skills := make(map[string]json.RawMessage)
	if err := json.Unmarshal(skillsRaw, &skills); err != nil {
		return Document{}, fmt.Errorf("%w: skills: %w", ErrInvalidDocument, err)
	}

	return Document{
		fields: fields,
		skills: skills,
	}, nil
}

func (d Document) Skill(name string) (SkillEntry, bool) {
	raw, ok := d.skills[name]
	if !ok {
		return SkillEntry{}, false
	}

	var decoded struct {
		SourceURL  string `json:"sourceUrl"`
		Source     string `json:"source"`
		SkillPath  string `json:"skillPath"`
		PluginName string `json:"pluginName"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return SkillEntry{}, false
	}

	return SkillEntry{
		SourceURL:  decoded.SourceURL,
		Source:     decoded.Source,
		SkillPath:  decoded.SkillPath,
		PluginName: decoded.PluginName,
		Raw:        append(json.RawMessage(nil), raw...),
	}, true
}

func (d Document) SkillNames() []string {
	return sortedNames(d.skills)
}

func (d Document) RemoveSkill(name string) (bool, error) {
	if !validSkillName(name) {
		return false, ErrInvalidSkillName
	}
	if _, ok := d.skills[name]; !ok {
		return false, nil
	}
	delete(d.skills, name)
	return true, nil
}

func (d Document) Marshal() ([]byte, error) {
	fields := cloneFields(d.fields)
	skillsRaw, err := marshalSkills(d.skills)
	if err != nil {
		return nil, err
	}
	fields["skills"] = skillsRaw

	raw, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("marshal lockfile: %w", err)
	}
	return append(raw, '\n'), nil
}

func cloneFields(fields map[string]json.RawMessage) map[string]json.RawMessage {
	clone := make(map[string]json.RawMessage, len(fields))
	for key, value := range fields {
		clone[key] = append(json.RawMessage(nil), value...)
	}
	return clone
}

func marshalSkills(skills map[string]json.RawMessage) (json.RawMessage, error) {
	ordered := make(map[string]json.RawMessage, len(skills))
	for _, name := range sortedNames(skills) {
		ordered[name] = append(json.RawMessage(nil), skills[name]...)
	}
	raw, err := json.Marshal(ordered)
	if err != nil {
		return nil, fmt.Errorf("marshal skills: %w", err)
	}
	return raw, nil
}

func sortedNames(skills map[string]json.RawMessage) []string {
	names := make([]string, 0, len(skills))
	for name := range skills {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func validSkillName(name string) bool {
	return name != "" && name != "." && name != ".." && filepath.Base(name) == name
}
