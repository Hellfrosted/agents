package output

import (
	"encoding/json"
	"fmt"
	"io"
)

func WriteSummaryJSON(writer io.Writer, summary Summary) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(summary); err != nil {
		return fmt.Errorf("write summary json: %w", err)
	}
	return nil
}

func WriteJSONLEvent(writer io.Writer, event Event) error {
	encoder := json.NewEncoder(writer)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(event); err != nil {
		return fmt.Errorf("write jsonl event: %w", err)
	}
	return nil
}
