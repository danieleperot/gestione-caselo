package main

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/daniele/gestione-caselo/backend/emails/internal/emails/admin_new_event"
	"github.com/daniele/gestione-caselo/backend/emails/internal/emails/customer_new_event"
	"github.com/daniele/gestione-caselo/backend/emails/internal/handler"
	"github.com/daniele/gestione-caselo/backend/emails/internal/mailer"
	"github.com/daniele/gestione-caselo/backend/emails/internal/message"
	"github.com/daniele/gestione-caselo/backend/emails/internal/registry"
	appconfig "github.com/daniele/gestione-caselo/backend/internal/config"
)

var h *handler.Handler

func init() {
	ctx := context.Background()

	adminEmailsStr := appconfig.GetEnvVariable("ADMIN_EMAILS")
	if adminEmailsStr == "" {
		log.Fatal("ADMIN_EMAILS is required")
	}
	adminEmails := strings.Split(adminEmailsStr, ",")
	for i, email := range adminEmails {
		adminEmails[i] = strings.TrimSpace(email)
	}

	fromAddr := appconfig.GetEnvVariable("FROM_ADDRESS")
	if fromAddr == "" {
		fromAddr = "noreply@gestione-caselo.it"
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	sesClient := ses.NewFromConfig(cfg)
	m := mailer.New(sesClient)

	reg := registry.New()
	reg.Register("admin_new_event", func(metadata map[string]interface{}) (message.EmailMessage, error) {
		return admin_new_event.New(metadata, adminEmails)
	})
	reg.Register("customer_new_event", func(metadata map[string]interface{}) (message.EmailMessage, error) {
		return customer_new_event.New(metadata)
	})

	h = handler.New(reg, m, fromAddr)
}

func handleRequest(ctx context.Context, sqsEvent events.SQSEvent) error {
	return h.ProcessSQSEvent(ctx, sqsEvent)
}

func main() {
	lambda.Start(handleRequest)
}
