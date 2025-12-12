package handlers

import (
	"context"
	"fmt"
	"log"

	"firebase.google.com/go/v4/messaging"
	"github.com/lwinmoehein/cec-notifications/config"
)

// SendFCMNotification sends a Firebase Cloud Messaging notification
func SendFCMNotification(ctx context.Context, msg *NotificationMessage) error {
	client, err := config.GetMessagingClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FCM client: %w", err)
	}

	// Build the FCM message
	message := &messaging.Message{
		Token: msg.FCMToken,
		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},
		Data: msg.Data,
	}

	// Send the message
	response, err := client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send FCM message: %w", err)
	}

	log.Printf("Successfully sent FCM message: %s", response)
	return nil
}

// SendFCMNotificationBatch sends multiple FCM notifications
func SendFCMNotificationBatch(ctx context.Context, messages []*NotificationMessage) ([]error, error) {
	client, err := config.GetMessagingClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get FCM client: %w", err)
	}

	// Build FCM messages
	fcmMessages := make([]*messaging.Message, len(messages))
	for i, msg := range messages {
		fcmMessages[i] = &messaging.Message{
			Token: msg.FCMToken,
			Notification: &messaging.Notification{
				Title: msg.Title,
				Body:  msg.Body,
			},
			Data: msg.Data,
		}
	}

	// Send messages as a batch
	response, err := client.SendAll(ctx, fcmMessages)
	if err != nil {
		return nil, fmt.Errorf("failed to send FCM messages: %w", err)
	}

	log.Printf("Successfully sent %d/%d FCM messages", response.SuccessCount, len(messages))

	// Collect errors from failed messages
	var errors []error
	for i, resp := range response.Responses {
		if !resp.Success {
			errors = append(errors, fmt.Errorf("message %d failed: %v", i, resp.Error))
		}
	}

	return errors, nil
}
