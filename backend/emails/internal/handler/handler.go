package handler

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
)

type Registry interface {
	Get(sqsMessageBody string) (message.EmailMessage, error)
}

type Mailer interface {
	SendEmail(ctx context.Context, from string, to []string, subject, body string, replyTo []string) error
}

type Handler struct {
	registry *Registry
	mailer   Mailer
	fromAddr string
}

func New(registry Registry, mailer Mailer, fromAddr string) *Handler {
	return &Handler{
		registry: &registry,
		mailer:   mailer,
		fromAddr: fromAddr,
	}
}

func (h *Handler) ProcessSQSEvent(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		if err := h.processMessage(ctx, record); err != nil {
			log.Printf("Error processing message %s: %v", record.MessageId, err)
			return err
		}
	}
	return nil
}

func (h *Handler) processMessage(ctx context.Context, record events.SQSMessage) error {
	emailMsg, err := (*h.registry).Get(record.Body)
	if err != nil {
		return fmt.Errorf("failed to load email from registry: %w", err)
	}

	recipients, err := emailMsg.Recipients()
	if err != nil {
		return fmt.Errorf("failed to get recipients: %w", err)
	}

	body, err := emailMsg.Render()
	if err != nil {
		return fmt.Errorf("failed to render email: %w", err)
	}

	subject := emailMsg.Subject()
	replyTo := emailMsg.ReplyTo()

	if err := h.mailer.SendEmail(ctx, h.fromAddr, recipients, subject, body, replyTo); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("Successfully sent email for message %s", record.MessageId)
	return nil
}
