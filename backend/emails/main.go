package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// handler processes SQS messages and sends emails via SES
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	sesClient := ses.NewFromConfig(cfg)

	for _, message := range sqsEvent.Records {
		fmt.Printf("Processing message ID: %s\n", message.MessageId)
		fmt.Printf("Message body: %s\n", message.Body)

		// TODO: Parse message body into email struct
		// TODO: Validate email parameters
		// TODO: Call sesClient.SendEmail() with proper parameters
		// TODO: Handle SES errors and implement retry logic

		_ = sesClient // Suppress unused variable warning for now
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
