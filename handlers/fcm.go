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

func SubscribeTokenToTopic(ctx context.Context, msg *NotificationMessage) error {

	return nil
}

func UnsubscribeTokenFromTopic(ctx context.Context, msg *NotificationMessage) error {

	return nil
}

func SendFCMNotificationToTopic(ctx context.Context, msg *NotificationMessage) error {

	return nil
}
