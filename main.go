package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lwinmoehein/cec-notifications/handlers"
)

// Handler processes SQS events and sends FCM notifications
func Handler(ctx context.Context, sqsEvent events.SQSEvent) (events.SQSEventResponse, error) {
	log.Printf("Processing %d SQS messages", len(sqsEvent.Records))

	var batchItemFailures []events.SQSBatchItemFailure

	for _, record := range sqsEvent.Records {
		if err := processMessage(ctx, record); err != nil {
			log.Printf("Error processing message %s: %v", record.MessageId, err)
			// Add to batch failures for retry
			batchItemFailures = append(batchItemFailures, events.SQSBatchItemFailure{
				ItemIdentifier: record.MessageId,
			})
		}
	}

	// Return batch item failures for SQS to retry
	return events.SQSEventResponse{
		BatchItemFailures: batchItemFailures,
	}, nil
}

func processMessage(ctx context.Context, record events.SQSMessage) error {
	// Parse the SQS message
	msg, err := handlers.ParseSQSMessage(record.Body)
	if err != nil {
		return fmt.Errorf("failed to parse SQS message: %w", err)
	}

	log.Printf("Sending FCM notification to token: %s, title: %s", msg.FCMToken, msg.Title)

	// Send FCM notification
	if err := handlers.SendFCMNotification(ctx, msg); err != nil {
		return fmt.Errorf("failed to send FCM notification: %w", err)
	}

	log.Printf("Successfully processed message %s", record.MessageId)
	return nil
}

func main() {
	lambda.Start(Handler)
}
