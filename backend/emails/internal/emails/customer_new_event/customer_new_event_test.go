package customer_new_event

import (
	"strings"
	"testing"
)

func TestNew_ValidMetadata(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@test.com",
		"eventDate":     "2025-12-01",
		"eventType":     "Wedding",
	}

	email, err := New(metadata)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if email == nil {
		t.Fatal("Expected email, got nil")
	}
}

func TestNew_MissingEventId(t *testing.T) {
	metadata := map[string]interface{}{
		"customerEmail": "customer@test.com",
	}

	_, err := New(metadata)
	if err == nil {
		t.Error("Expected error for missing eventId")
	}
}

func TestNew_MissingCustomerEmail(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId": "evt_123",
	}

	_, err := New(metadata)
	if err == nil {
		t.Error("Expected error for missing customerEmail")
	}
}

func TestRecipients(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@test.com",
	}

	email, _ := New(metadata)

	recipients, err := email.Recipients()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(recipients) != 1 {
		t.Errorf("Expected 1 recipient, got %d", len(recipients))
	}

	if recipients[0] != "customer@test.com" {
		t.Errorf("Recipient = %v, want customer@test.com", recipients[0])
	}
}

func TestSubject(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@test.com",
	}

	email, _ := New(metadata)

	subject := email.Subject()
	expected := "Event Submission Confirmed - evt_123"

	if subject != expected {
		t.Errorf("Subject = %v, want %v", subject, expected)
	}
}

func TestRender(t *testing.T) {
	metadata := map[string]interface{}{
		"eventId":       "evt_123",
		"customerEmail": "customer@test.com",
		"eventDate":     "2025-12-01",
		"eventType":     "Wedding",
	}

	email, _ := New(metadata)

	body, err := email.Render()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if !strings.Contains(body, "evt_123") {
		t.Error("Body should contain eventId")
	}

	if !strings.Contains(body, "2025-12-01") {
		t.Error("Body should contain eventDate")
	}

	if !strings.Contains(body, "Wedding") {
		t.Error("Body should contain eventType")
	}
}
