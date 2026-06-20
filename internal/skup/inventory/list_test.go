package inventory

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestListInstalled_returnsSortedSkillDirectories_whenSkillsExist(t *testing.T) {
	// Given
	agentsHome := t.TempDir()
	skillsDir := filepath.Join(agentsHome, "skills")
	for _, name := range []string{"zeta", "alpha"} {
		if err := os.MkdirAll(filepath.Join(skillsDir, name), 0o700); err != nil {
			t.Fatalf("MkdirAll returned error: %v", err)
		}
	}
	if err := os.WriteFile(filepath.Join(skillsDir, "README.md"), []byte("ignore"), 0o600); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}

	// When
	got, err := ListInstalled(agentsHome)

	// Then
	if err != nil {
		t.Fatalf("ListInstalled returned error: %v", err)
	}
	want := []string{"alpha", "zeta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("skills = %#v, want %#v", got, want)
	}
}

func TestListInstalled_returnsEmptyList_whenSkillsDirectoryMissing(t *testing.T) {
	// Given
	agentsHome := t.TempDir()

	// When
	got, err := ListInstalled(agentsHome)

	// Then
	if err != nil {
		t.Fatalf("ListInstalled returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("skills = %#v, want empty", got)
	}
}
