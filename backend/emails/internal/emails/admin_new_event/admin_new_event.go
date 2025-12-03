package admin_new_event

import (
	_ "embed"
	"fmt"

	"github.com/daniele/gestione-caselo/backend/emails/internal/emails"
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
	return "Nuova Richiesta Evento"
}

func (e *AdminNewEvent) ReplyTo() []string {
	customerEmail := e.metadata["customerEmail"].(string)
	return []string{customerEmail}
}

func (e *AdminNewEvent) Render() (string, error) {
	return emails.RenderTemplate("admin_new_event", templateContent, e.metadata)
}
