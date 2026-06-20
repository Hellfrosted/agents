package status

import "github.com/Hellfrosted/agents/internal/skup/output"

func emitProgress(progress ProgressFunc, event output.Event) {
	if progress == nil {
		return
	}
	progress(event)
}

func repoProgressEvent(event output.EventName, source skillSource) output.Event {
	return output.Event{
		Type:      output.EventTypeEvent,
		Event:     event,
		SourceURL: source.sourceURL,
	}
}

func skillProgressEvent(event output.EventName, source skillSource) output.Event {
	return output.Event{
		Type:      output.EventTypeEvent,
		Event:     event,
		Name:      source.name,
		SourceURL: source.sourceURL,
	}
}
