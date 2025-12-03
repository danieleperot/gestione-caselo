package handler

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

type mockRegistry struct {
	getFunc func(sqsMessageBody string) (message.EmailMessage, error)
}

func (m *mockRegistry) Get(sqsMessageBody string) (message.EmailMessage, error) {
	return m.getFunc(sqsMessageBody)
}

type mockMailer struct {
	sendEmailFunc func(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error
}

func (m *mockMailer) SendEmail(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error {
	return m.sendEmailFunc(ctx, from, to, subject, body, replyTo)
}

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

func (m *mockEmailMessage) ReplyTo() []string {
	return nil
}

func TestProcessSQSEvent_Success(t *testing.T) {
	var capturedFrom string
	var capturedTo []string
	var capturedSubject string
	var capturedBody string

	mockReg := &mockRegistry{
		getFunc: func(sqsMessageBody string) (message.EmailMessage, error) {
			return &mockEmailMessage{
				recipients: []string{"test@example.com"},
				subject:    "Test Subject",
				body:       "Test Body",
			}, nil
		},
	}

	mockMail := &mockMailer{
		sendEmailFunc: func(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error {
			capturedFrom = from
			capturedTo = to
			capturedSubject = subject
			capturedBody = body
			return nil
		},
	}

	h := New(mockReg, mockMail, "from@test.com")

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId: "test-msg-1",
				Body:      `{"template":"test","metadata":{}}`,
			},
		},
	}

	err := h.ProcessSQSEvent(context.Background(), sqsEvent)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if capturedFrom != "from@test.com" {
		t.Errorf("From = %v, want from@test.com", capturedFrom)
	}

	if len(capturedTo) != 1 || capturedTo[0] != "test@example.com" {
		t.Errorf("To = %v, want [test@example.com]", capturedTo)
	}

	if capturedSubject != "Test Subject" {
		t.Errorf("Subject = %v, want Test Subject", capturedSubject)
	}

	if capturedBody != "Test Body" {
		t.Errorf("Body = %v, want Test Body", capturedBody)
	}
}

func TestProcessSQSEvent_MultipleMessages(t *testing.T) {
	processedCount := 0

	mockReg := &mockRegistry{
		getFunc: func(sqsMessageBody string) (message.EmailMessage, error) {
			return &mockEmailMessage{
				recipients: []string{"test@example.com"},
				subject:    "Test",
				body:       "Test",
			}, nil
		},
	}

	mockMail := &mockMailer{
		sendEmailFunc: func(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error {
			processedCount++
			return nil
		},
	}

	h := New(mockReg, mockMail, "from@test.com")

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "msg-1", Body: `{"template":"test","metadata":{}}`},
			{MessageId: "msg-2", Body: `{"template":"test","metadata":{}}`},
			{MessageId: "msg-3", Body: `{"template":"test","metadata":{}}`},
		},
	}

	err := h.ProcessSQSEvent(context.Background(), sqsEvent)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if processedCount != 3 {
		t.Errorf("Processed %d messages, want 3", processedCount)
	}
}

func TestProcessSQSEvent_RegistryError(t *testing.T) {
	mockReg := &mockRegistry{
		getFunc: func(sqsMessageBody string) (message.EmailMessage, error) {
			return nil, fmt.Errorf("unknown template")
		},
	}

	mockMail := &mockMailer{
		sendEmailFunc: func(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error {
			return nil
		},
	}

	h := New(mockReg, mockMail, "from@test.com")

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "test-msg-1", Body: `{}`},
		},
	}

	err := h.ProcessSQSEvent(context.Background(), sqsEvent)
	if err == nil {
		t.Error("Expected error to trigger retry/DLQ flow, got nil")
	}
}

func TestProcessSQSEvent_MailerError(t *testing.T) {
	mockReg := &mockRegistry{
		getFunc: func(sqsMessageBody string) (message.EmailMessage, error) {
			return &mockEmailMessage{
				recipients: []string{"test@example.com"},
				subject:    "Test",
				body:       "Test",
			}, nil
		},
	}

	mockMail := &mockMailer{
		sendEmailFunc: func(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error {
			return fmt.Errorf("SES error")
		},
	}

	h := New(mockReg, mockMail, "from@test.com")

	sqsEvent := events.SQSEvent{
		Records: []events.SQSMessage{
			{MessageId: "test-msg-1", Body: `{"template":"test","metadata":{}}`},
		},
	}

	err := h.ProcessSQSEvent(context.Background(), sqsEvent)
	if err == nil {
		t.Error("Expected error to trigger retry/DLQ flow, got nil")
	}
}
