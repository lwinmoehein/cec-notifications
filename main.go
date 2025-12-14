package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/lwinmoehein/cec-notifications/handlers"
)

const (
	ActionSendSingle       = "SEND_SINGLE_NOTIFICATION"
	ActionSubscribeTopic   = "SUBSCRIBE_TO_TOPIC"
	ActionSendToTopic      = "SEND_TO_TOPIC_NOTIFICATION"
	ActionUnsubscribeTopic = "UNSUBSCRIBE_FROM_TOPIC"
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

	log.Printf("Processing message with ActionType: %s", msg.ActionType)
	switch msg.ActionType {
	case ActionSendSingle:
		if err := handlers.SendFCMNotification(ctx, msg); err != nil {
			return fmt.Errorf("failed to send FCM notification: %w", err)
		}

	case ActionSubscribeTopic:
		if err := handlers.SubscribeTokenToTopic(ctx, msg); err != nil {
			log.Printf("ERROR: Failed to subscribe token: %v", err)
			return err
		}

	case ActionSendToTopic:
		if err := handlers.SendFCMNotificationToTopic(ctx, msg); err != nil {
			log.Printf("ERROR: Failed to send topic notification: %v", err)
			return err
		}

	case ActionUnsubscribeTopic:
		if err := handlers.UnsubscribeTokenFromTopic(ctx, msg); err != nil {
			log.Printf("ERROR: Failed to unsubscribe token: %v", err)
			return err
		}

	default:
		log.Printf("WARNING: Unknown ActionType received: %s. Skipping message.", msg.ActionType)
		// No need to return an error for unknown types, let SQS delete it
	}

	log.Printf("Successfully processed message %s", record.MessageId)
	return nil
}

func main() {
	lambda.Start(Handler)
}
