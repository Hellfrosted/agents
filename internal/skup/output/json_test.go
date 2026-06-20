package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteSummaryJSON_emitsFinalDocument_whenSummaryProvided(t *testing.T) {
	// Given
	var buffer bytes.Buffer
	summary := Summary{
		OK:         true,
		Command:    "status",
		Entrypoint: "sk-up",
		DryRun:     false,
		Statuses: []SkillStatus{
			{
				Name:       "confidence-loop",
				Status:     StatusUpdate,
				SourceURL:  "https://github.com/example/skills.git",
				RemoteHash: "abc123",
			},
		},
	}

	// When
	err := WriteSummaryJSON(&buffer, summary)

	// Then
	if err != nil {
		t.Fatalf("WriteSummaryJSON returned error: %v", err)
	}
	var got map[string]json.RawMessage
	if err := json.Unmarshal(buffer.Bytes(), &got); err != nil {
		t.Fatalf("summary is not valid JSON: %v\n%s", err, buffer.String())
	}
	if !strings.HasSuffix(buffer.String(), "\n") {
		t.Fatalf("summary JSON must end with newline: %q", buffer.String())
	}
	if string(got["command"]) != `"status"` {
		t.Fatalf("command = %s, want status", got["command"])
	}
	if _, ok := got["statuses"]; !ok {
		t.Fatalf("summary missing statuses: %s", buffer.String())
	}
}

func TestWriteJSONLEvent_emitsSingleLineEvent_whenEventProvided(t *testing.T) {
	// Given
	var buffer bytes.Buffer
	event := Event{
		Type:      EventTypeEvent,
		Event:     EventFetch,
		SourceURL: "https://github.com/example/skills.git",
	}

	// When
	err := WriteJSONLEvent(&buffer, event)

	// Then
	if err != nil {
		t.Fatalf("WriteJSONLEvent returned error: %v", err)
	}
	lines := strings.Split(strings.TrimSuffix(buffer.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("line count = %d, want 1 in %q", len(lines), buffer.String())
	}
	var got Event
	if err := json.Unmarshal([]byte(lines[0]), &got); err != nil {
		t.Fatalf("event is not valid JSON: %v\n%s", err, lines[0])
	}
	if got.Event != EventFetch {
		t.Fatalf("event = %q, want %q", got.Event, EventFetch)
	}
}
