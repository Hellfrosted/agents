package lockfile

import (
	"encoding/json"
	"errors"
	"testing"
)

const sampleLock = `{
  "version": 1,
  "workspace": {"owner": "local"},
  "skills": {
    "confidence-loop": {
      "sourceUrl": "https://github.com/example/skills.git",
      "skillPath": "skills/confidence-loop/SKILL.md",
      "pluginName": "example-skills"
    },
    "legacy": {
      "source": "owner/repo",
      "skillPath": "skills/legacy/SKILL.md",
      "extra": {"keep": true}
    }
  }
}`

func TestParse_readsSkillEntries_whenLockfileContainsSourceForms(t *testing.T) {
	// Given
	raw := []byte(sampleLock)

	// When
	doc, err := Parse(raw)

	// Then
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	entry, ok := doc.Skill("confidence-loop")
	if !ok {
		t.Fatal("confidence-loop missing")
	}
	if entry.SourceURL != "https://github.com/example/skills.git" {
		t.Fatalf("SourceURL = %q", entry.SourceURL)
	}
	legacy, ok := doc.Skill("legacy")
	if !ok {
		t.Fatal("legacy missing")
	}
	if legacy.Source != "owner/repo" {
		t.Fatalf("Source = %q", legacy.Source)
	}
}

func TestRemoveSkill_preservesUnrelatedEntriesAndTopLevelFields(t *testing.T) {
	// Given
	doc, err := Parse([]byte(sampleLock))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// When
	removed, err := doc.RemoveSkill("confidence-loop")
	if err != nil {
		t.Fatalf("RemoveSkill returned error: %v", err)
	}
	raw, err := doc.Marshal()
	if err != nil {
		t.Fatalf("Marshal returned error: %v", err)
	}

	// Then
	if !removed {
		t.Fatal("removed = false, want true")
	}
	var decoded map[string]json.RawMessage
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("marshaled document is invalid JSON: %v", err)
	}
	if _, ok := decoded["workspace"]; !ok {
		t.Fatalf("workspace field was not preserved: %s", raw)
	}
	if _, ok := doc.Skill("confidence-loop"); ok {
		t.Fatal("confidence-loop still present")
	}
	if _, ok := doc.Skill("legacy"); !ok {
		t.Fatal("legacy skill was not preserved")
	}
}

func TestRemoveSkill_rejectsUnsafeSkillName(t *testing.T) {
	// Given
	doc, err := Parse([]byte(sampleLock))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	// When
	_, err = doc.RemoveSkill("../bad")

	// Then
	if !errors.Is(err, ErrInvalidSkillName) {
		t.Fatalf("RemoveSkill error = %v, want ErrInvalidSkillName", err)
	}
}
