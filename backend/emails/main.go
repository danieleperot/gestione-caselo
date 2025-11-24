package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// handler processes SQS messages and sends emails via SES
func handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, message := range sqsEvent.Records {
		fmt.Printf("Processing message ID: %s\n", message.MessageId)
		fmt.Printf("Message body: %s\n", message.Body)

		// TODO: Parse message body
		// TODO: Send email via SES
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
