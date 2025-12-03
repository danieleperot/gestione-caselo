package mailer

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type mockSESClient struct {
	sendEmailFunc func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

func (m *mockSESClient) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	return m.sendEmailFunc(ctx, params, optFns...)
}

func TestSendEmail_Success(t *testing.T) {
	mockClient := &mockSESClient{
		sendEmailFunc: func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
			return &ses.SendEmailOutput{
				MessageId: stringPtr("test-message-id"),
			}, nil
		},
	}

	m := &Mailer{client: mockClient}

	err := m.SendEmail(
		context.Background(),
		"from@test.com",
		[]string{"to@test.com"},
		"Test Subject",
		"Test Body",
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestSendEmail_MultipleRecipients(t *testing.T) {
	var capturedRecipients []string

	mockClient := &mockSESClient{
		sendEmailFunc: func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
			capturedRecipients = params.Destination.ToAddresses
			return &ses.SendEmailOutput{
				MessageId: stringPtr("test-message-id"),
			}, nil
		},
	}

	m := &Mailer{client: mockClient}

	recipients := []string{"to1@test.com", "to2@test.com", "to3@test.com"}
	err := m.SendEmail(
		context.Background(),
		"from@test.com",
		recipients,
		"Test Subject",
		"Test Body",
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(capturedRecipients) != 3 {
		t.Errorf("Expected 3 recipients, got %d", len(capturedRecipients))
	}
}

func TestSendEmail_SESError(t *testing.T) {
	mockClient := &mockSESClient{
		sendEmailFunc: func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
			return nil, fmt.Errorf("SES error")
		},
	}

	m := &Mailer{client: mockClient}

	err := m.SendEmail(
		context.Background(),
		"from@test.com",
		[]string{"to@test.com"},
		"Test Subject",
		"Test Body",
		nil,
	)

	if err == nil {
		t.Error("Expected error from SES, got nil")
	}
}

func stringPtr(s string) *string {
	return &s
}
