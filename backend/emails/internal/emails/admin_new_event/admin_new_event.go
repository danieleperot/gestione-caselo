package admin_new_event

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

//go:embed template.txt
var templateContent string

type AdminNewEvent struct {
	metadata    map[string]interface{}
	adminEmails []string
}

func New(metadata map[string]interface{}, adminEmails []string) (message.EmailMessage, error) {
	eventId, ok := metadata["eventId"].(string)
	if !ok || eventId == "" {
		return nil, fmt.Errorf("eventId is required in metadata")
	}

	if len(adminEmails) == 0 {
		return nil, fmt.Errorf("at least one admin email is required")
	}

	return &AdminNewEvent{
		metadata:    metadata,
		adminEmails: adminEmails,
	}, nil
}

func (e *AdminNewEvent) Recipients() ([]string, error) {
	return e.adminEmails, nil
}

func (e *AdminNewEvent) Subject() string {
	eventId := e.metadata["eventId"].(string)
	return fmt.Sprintf("New Event Request - %s", eventId)
}

func (e *AdminNewEvent) Render() (string, error) {
	tmpl, err := template.New("admin_new_event").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, e.metadata); err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return buf.String(), nil
}
