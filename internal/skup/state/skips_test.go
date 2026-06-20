package state

import (
	"errors"
	"reflect"
	"testing"
)

const sampleSkips = `{
  "skips": {
    "confidence-loop": {
      "remoteHash": "abc123",
      "sourceUrl": "https://github.com/example/skills.git",
      "skippedAt": "2026-06-20T00:00:00Z"
    }
  }
}`

func TestParseSkips_readsUpstreamHashSkip_whenStateExists(t *testing.T) {
	// Given
	raw := []byte(sampleSkips)

	// When
	doc, err := ParseSkips(raw)

	// Then
	if err != nil {
		t.Fatalf("ParseSkips returned error: %v", err)
	}
	entry, ok := doc.Entry("confidence-loop")
	if !ok {
		t.Fatal("confidence-loop skip missing")
	}
	if entry.RemoteHash != "abc123" {
		t.Fatalf("RemoteHash = %q, want abc123", entry.RemoteHash)
	}
}

func TestSetEntry_writesUpstreamHashSkip_whenEntryProvided(t *testing.T) {
	// Given
	doc := NewSkips()
	entry := SkipEntry{
		RemoteHash: "def456",
		SourceURL:  "https://github.com/example/skills.git",
		SkippedAt:  "2026-06-20T01:00:00Z",
	}

	// When
	err := doc.SetEntry("confidence-loop", entry)

	// Then
	if err != nil {
		t.Fatalf("SetEntry returned error: %v", err)
	}
	got, ok := doc.Entry("confidence-loop")
	if !ok {
		t.Fatal("confidence-loop skip missing")
	}
	if got.RemoteHash != entry.RemoteHash {
		t.Fatalf("RemoteHash = %q, want %q", got.RemoteHash, entry.RemoteHash)
	}
}

func TestRemoveEntry_deletesOnlyRequestedSkip(t *testing.T) {
	// Given
	doc, err := ParseSkips([]byte(sampleSkips))
	if err != nil {
		t.Fatalf("ParseSkips returned error: %v", err)
	}

	// When
	removed, err := doc.RemoveEntry("confidence-loop")

	// Then
	if err != nil {
		t.Fatalf("RemoveEntry returned error: %v", err)
	}
	if !removed {
		t.Fatal("removed = false, want true")
	}
	if _, ok := doc.Entry("confidence-loop"); ok {
		t.Fatal("confidence-loop skip still present")
	}
}

func TestNames_returnsSortedSkillNames(t *testing.T) {
	// Given
	doc := NewSkips()
	if err := doc.SetEntry("zeta", SkipEntry{RemoteHash: "hash-z"}); err != nil {
		t.Fatalf("SetEntry zeta returned error: %v", err)
	}
	if err := doc.SetEntry("alpha", SkipEntry{RemoteHash: "hash-a"}); err != nil {
		t.Fatalf("SetEntry alpha returned error: %v", err)
	}

	// When
	got := doc.Names()

	// Then
	want := []string{"alpha", "zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Names = %#v, want %#v", got, want)
	}
}

func TestSetEntry_rejectsUnsafeSkillName(t *testing.T) {
	// Given
	doc := NewSkips()

	// When
	err := doc.SetEntry("../bad", SkipEntry{RemoteHash: "abc"})

	// Then
	if !errors.Is(err, ErrInvalidSkillName) {
		t.Fatalf("SetEntry error = %v, want ErrInvalidSkillName", err)
	}
}
