package customer_new_event

import (
	_ "embed"
	"fmt"

	"github.com/daniele/gestione-caselo/backend/emails/internal/emails"
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
	return "Richiesta Evento Confermata"
}

func (e *CustomerNewEvent) Render() (string, error) {
	return emails.RenderTemplate("customer_new_event", templateContent, e.metadata)
}
