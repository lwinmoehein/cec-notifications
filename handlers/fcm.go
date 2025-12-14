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

// SubscribeTokenToTopic subscribes a device token to an FCM topic
func SubscribeTokenToTopic(ctx context.Context, msg *NotificationMessage) error {
	// Validate required fields
	if msg.TopicName == "" {
		return fmt.Errorf("topicName is required for SUBSCRIBE_TO_TOPIC action")
	}
	if msg.FCMToken == "" {
		return fmt.Errorf("fcmToken is required for SUBSCRIBE_TO_TOPIC action")
	}

	client, err := config.GetMessagingClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FCM client: %w", err)
	}

	// Subscribe the token to the topic
	log.Printf("Subscribing token to topic: %s", msg.TopicName)
	response, err := client.SubscribeToTopic(ctx, []string{msg.FCMToken}, msg.TopicName)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", msg.TopicName, err)
	}

	// Check for failures
	if response.FailureCount > 0 {
		log.Printf("WARNING: %d tokens failed to subscribe to topic %s", response.FailureCount, msg.TopicName)
		for idx, err := range response.Errors {
			log.Printf("  Error [%d]: %v", idx, err)
		}
		return fmt.Errorf("failed to subscribe %d tokens to topic %s", response.FailureCount, msg.TopicName)
	}

	log.Printf("✓ Successfully subscribed token to topic: %s (Success: %d)", msg.TopicName, response.SuccessCount)
	return nil
}

// UnsubscribeTokenFromTopic unsubscribes a device token from an FCM topic
func UnsubscribeTokenFromTopic(ctx context.Context, msg *NotificationMessage) error {
	// Validate required fields
	if msg.TopicName == "" {
		return fmt.Errorf("topicName is required for UNSUBSCRIBE_FROM_TOPIC action")
	}
	if msg.FCMToken == "" {
		return fmt.Errorf("fcmToken is required for UNSUBSCRIBE_FROM_TOPIC action")
	}

	client, err := config.GetMessagingClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FCM client: %w", err)
	}

	// Unsubscribe the token from the topic
	log.Printf("Unsubscribing token from topic: %s", msg.TopicName)
	response, err := client.UnsubscribeFromTopic(ctx, []string{msg.FCMToken}, msg.TopicName)
	if err != nil {
		return fmt.Errorf("failed to unsubscribe from topic %s: %w", msg.TopicName, err)
	}

	// Check for failures
	if response.FailureCount > 0 {
		log.Printf("WARNING: %d tokens failed to unsubscribe from topic %s", response.FailureCount, msg.TopicName)
		for idx, err := range response.Errors {
			log.Printf("  Error [%d]: %v", idx, err)
		}
		return fmt.Errorf("failed to unsubscribe %d tokens from topic %s", response.FailureCount, msg.TopicName)
	}

	log.Printf("✓ Successfully unsubscribed token from topic: %s (Success: %d)", msg.TopicName, response.SuccessCount)
	return nil
}

// SendFCMNotificationToTopic sends a notification to all subscribers of a topic
func SendFCMNotificationToTopic(ctx context.Context, msg *NotificationMessage) error {
	// Validate required fields
	if msg.TopicName == "" {
		return fmt.Errorf("topicName is required for SEND_TO_TOPIC_NOTIFICATION action")
	}
	if msg.Title == "" {
		return fmt.Errorf("title is required for SEND_TO_TOPIC_NOTIFICATION action")
	}
	if msg.Body == "" {
		return fmt.Errorf("body is required for SEND_TO_TOPIC_NOTIFICATION action")
	}

	client, err := config.GetMessagingClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to get FCM client: %w", err)
	}

	// Build the FCM message for topic
	message := &messaging.Message{
		Topic: msg.TopicName,
		Notification: &messaging.Notification{
			Title: msg.Title,
			Body:  msg.Body,
		},
		Data: msg.Data,
	}

	// Send the message to the topic
	log.Printf("Sending notification to topic: %s (Title: %s)", msg.TopicName, msg.Title)
	response, err := client.Send(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send FCM message to topic %s: %w", msg.TopicName, err)
	}

	log.Printf("✓ Successfully sent FCM message to topic %s: %s", msg.TopicName, response)
	return nil
}
