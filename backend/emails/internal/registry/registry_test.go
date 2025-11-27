package registry

import (
	"fmt"
	"testing"

	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

type mockEmailMessage struct {
	recipients []string
	subject    string
	body       string
}

func (m *mockEmailMessage) Recipients() ([]string, error) {
	return m.recipients, nil
}

func (m *mockEmailMessage) Subject() string {
	return m.subject
}

func (m *mockEmailMessage) Render() (string, error) {
	return m.body, nil
}

func TestRegistry_Register(t *testing.T) {
	r := New()

	constructor := func(metadata map[string]interface{}) (message.EmailMessage, error) {
		return &mockEmailMessage{
			recipients: []string{"test@example.com"},
			subject:    "Test",
			body:       "Test body",
		}, nil
	}

	r.Register("test_template", constructor)

	if r.constructors["test_template"] == nil {
		t.Error("Constructor was not registered")
	}
}

func TestRegistry_Get_ValidMessage(t *testing.T) {
	r := New()

	r.Register("test_template", func(metadata map[string]interface{}) (message.EmailMessage, error) {
		return &mockEmailMessage{
			recipients: []string{"test@example.com"},
			subject:    "Test Subject",
			body:       "Test Body",
		}, nil
	})

	sqsMessage := `{
		"template": "test_template",
		"metadata": {
			"eventId": "123"
		}
	}`

	msg, err := r.Get(sqsMessage)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if msg == nil {
		t.Fatal("Expected message, got nil")
	}

	subject := msg.Subject()
	if subject != "Test Subject" {
		t.Errorf("Subject = %v, want Test Subject", subject)
	}
}

func TestRegistry_Get_UnknownTemplate(t *testing.T) {
	r := New()

	sqsMessage := `{
		"template": "unknown_template",
		"metadata": {}
	}`

	_, err := r.Get(sqsMessage)
	if err == nil {
		t.Error("Expected error for unknown template, got nil")
	}
}

func TestRegistry_Get_InvalidJSON(t *testing.T) {
	r := New()

	sqsMessage := `{invalid json}`

	_, err := r.Get(sqsMessage)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestRegistry_Get_MissingTemplate(t *testing.T) {
	r := New()

	sqsMessage := `{
		"metadata": {}
	}`

	_, err := r.Get(sqsMessage)
	if err == nil {
		t.Error("Expected error for missing template field, got nil")
	}
}

func TestRegistry_Get_ConstructorError(t *testing.T) {
	r := New()

	r.Register("error_template", func(metadata map[string]interface{}) (message.EmailMessage, error) {
		return nil, fmt.Errorf("constructor error")
	})

	sqsMessage := `{
		"template": "error_template",
		"metadata": {}
	}`

	_, err := r.Get(sqsMessage)
	if err == nil {
		t.Error("Expected constructor error to be propagated")
	}
}
