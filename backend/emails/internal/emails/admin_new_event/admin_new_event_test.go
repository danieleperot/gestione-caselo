package admin_new_event

import (
	"strings"
	"testing"
)

func TestNew_ValidMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@example.com",
		"eventDate":     "2025-12-01",
		"eventType":     "Wedding",
	}
	adminEmails := []string{"admin@test.com", "manager@test.com"}

	email, err := New(metadata, adminEmails)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if email == nil {
		t.Fatal("Expected email, got nil")
	}
}

func TestNew_MissingEventId(t *testing.T) {
	metadata := map[string]interface{}{
		"customerEmail": "customer@example.com",
	}
	adminEmails := []string{"admin@test.com"}

	_, err := New(metadata, adminEmails)
	if err == nil {
		t.Error("Expected error for missing eventId")
	}
}

func TestNew_EmptyAdminEmails(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId": "evt_123",
	}
	adminEmails := []string{}

	_, err := New(metadata, adminEmails)
	if err == nil {
		t.Error("Expected error for empty admin emails")
	}
}

func TestRecipients(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId": "evt_123",
	}
	adminEmails := []string{"admin@test.com", "manager@test.com"}

	email, _ := New(metadata, adminEmails)

	recipients, err := email.Recipients()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(recipients) != 2 {
		t.Errorf("Expected 2 recipients, got %d", len(recipients))
	}

	if recipients[0] != "admin@test.com" {
		t.Errorf("First recipient = %v, want admin@test.com", recipients[0])
	}

	if recipients[1] != "manager@test.com" {
		t.Errorf("Second recipient = %v, want manager@test.com", recipients[1])
	}
}

func TestSubject(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId": "evt_123",
	}
	adminEmails := []string{"admin@test.com"}

	email, _ := New(metadata, adminEmails)

	subject := email.Subject()
	expected := "Nuova Richiesta Evento"

	if subject != expected {
		t.Errorf("Subject = %v, want %v", subject, expected)
	}
}

func TestRender(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@example.com",
		"eventDate":     "2025-12-01",
		"eventType":     "Wedding",
	}
	adminEmails := []string{"admin@test.com"}

	email, _ := New(metadata, adminEmails)

	body, err := email.Render()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(body, "evt_123") {
		t.Error("Body should contain eventId")
	}

	if !strings.Contains(body, "customer@example.com") {
		t.Error("Body should contain customerEmail")
	}

	if !strings.Contains(body, "2025-12-01") {
		t.Error("Body should contain eventDate")
	}

	if !strings.Contains(body, "Wedding") {
		t.Error("Body should contain eventType")
	}
}
