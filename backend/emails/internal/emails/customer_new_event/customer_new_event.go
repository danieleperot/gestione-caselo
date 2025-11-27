package customer_new_event

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

//go:embed template.txt
var templateContent string

type CustomerNewEvent struct {
	metadata map[string]interface{}
}

func New(metadata map[string]interface{}) (message.EmailMessage, error) {
	eventId, ok := metadata["eventId"].(string)
	if !ok || eventId == "" {
		return nil, fmt.Errorf("eventId is required in metadata")
	}

	customerEmail, ok := metadata["customerEmail"].(string)
	if !ok || customerEmail == "" {
		return nil, fmt.Errorf("customerEmail is required in metadata")
	}

	return &CustomerNewEvent{
		metadata: metadata,
	}, nil
}

func (e *CustomerNewEvent) Recipients() ([]string, error) {
	customerEmail := e.metadata["customerEmail"].(string)
	return []string{customerEmail}, nil
}

func (e *CustomerNewEvent) Subject() string {
	eventId := e.metadata["eventId"].(string)
	return fmt.Sprintf("Event Submission Confirmed - %s", eventId)
}

func (e *CustomerNewEvent) Render() (string, error) {
	tmpl, err := template.New("customer_new_event").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, e.metadata); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
