package app

import (
	"fmt"
	"io"
	"sync"

	"github.com/Hellfrosted/agents/internal/skup/output"
	"github.com/Hellfrosted/agents/internal/skup/status"
)

func humanProgress(writer io.Writer, colorMode string) status.ProgressFunc {
	if writer == nil {
		return nil
	}
	var mu sync.Mutex
	return func(event output.Event) {
		line := humanProgressLine(event, colorMode, writer)
		if line == "" {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		writeText(writer, line)
	}
}

func humanProgressLine(event output.Event, colorMode string, writer io.Writer) string {
	label, ok := humanProgressLabel(event.Event, colorMode, writer)
	if !ok {
		return ""
	}
	target := event.Name
	if target == "" {
		target = event.SourceURL
	}
	if target == "" {
		return ""
	}
	return fmt.Sprintf("%s %s\n", label, target)
}

func humanProgressLabel(event output.EventName, colorMode string, writer io.Writer) (string, bool) {
	text, ok := progressLabelText(event)
	if !ok {
		return "", false
	}
	label := fmt.Sprintf("%-7s", text)
	if !statusColorEnabled(colorMode, writer) {
		return label, true
	}
	return progressColor(event) + label + "\x1b[0m", true
}

func progressLabelText(event output.EventName) (string, bool) {
	switch event {
	case output.EventFetch:
		return "FETCH", true
	case output.EventClone:
		return "CLONE", true
	case output.EventCompare:
		return "CHECK", true
	case output.EventRepair:
		return "REPAIR", true
	default:
		return "", false
	}
}

func progressColor(event output.EventName) string {
	switch event {
	case output.EventRepair:
		return "\x1b[33m"
	case output.EventCompare:
		return "\x1b[36m"
	default:
		return "\x1b[34m"
	}
}
