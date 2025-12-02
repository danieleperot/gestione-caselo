package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	queueURL := os.Getenv("SQS_QUEUE_URL")
	if queueURL == "" {
		queueURL = "http://elasticmq:9324/000000000000/emails"
	}

	lambdaURL := os.Getenv("LAMBDA_URL")
	if lambdaURL == "" {
		lambdaURL = "http://emails:8080/2015-03-31/functions/function/invocations"
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	sqsClient := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String("http://elasticmq:9324")
	})

	log.Printf("Starting SQS poller")
	log.Printf("Queue URL: %s", queueURL)
	log.Printf("Lambda URL: %s", lambdaURL)

	for {
		result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queueURL),
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		})

		if err != nil {
			log.Printf("Error receiving messages: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if len(result.Messages) == 0 {
			continue
		}

		for _, msg := range result.Messages {
			log.Printf("Received message: %s", *msg.MessageId)

			sqsEvent := events.SQSEvent{
				Records: []events.SQSMessage{
					{
						MessageId: *msg.MessageId,
						Body:      *msg.Body,
					},
				},
			}

			eventJSON, err := json.Marshal(sqsEvent)
			if err != nil {
				log.Printf("Error marshaling event: %v", err)
				continue
			}

			resp, err := http.Post(lambdaURL, "application/json", bytes.NewReader(eventJSON))
			if err != nil {
				log.Printf("Error invoking Lambda: %v", err)
				continue
			}

			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode != 200 {
				log.Printf("Lambda returned error (status %d): %s", resp.StatusCode, string(body))
				continue
			}

			log.Printf("Lambda processed successfully")

			_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				log.Printf("Error deleting message: %v", err)
			} else {
				log.Printf("Message deleted: %s", *msg.MessageId)
			}
		}
	}
}
