package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	sqsClient *sqs.Client
	queueURL  string
}

func New(sqsClient *sqs.Client, queueURL string) *Client {
	return &Client{
		sqsClient: sqsClient,
		queueURL:  queueURL,
	}
}

type EmailMessage struct {
	Template string                 `json:"template"`
	Metadata map[string]interface{} `json:"metadata"`
}

func (c *Client) SendEmailMessage(ctx context.Context, msg EmailMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	_, err = c.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("failed to send SQS message: %w", err)
	}

	return nil
}
