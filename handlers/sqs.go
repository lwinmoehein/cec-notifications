package handlers

import (
	"encoding/json"
	"fmt"
)

// NotificationMessage represents the expected SQS message format
type NotificationMessage struct {
	FCMToken string            `json:"fcmToken"`
	Title    string            `json:"title"`
	Body     string            `json:"body"`
	Data     map[string]string `json:"data,omitempty"`
}

// ParseSQSMessage parses the SQS message body into a NotificationMessage
func ParseSQSMessage(messageBody string) (*NotificationMessage, error) {
	var msg NotificationMessage
	if err := json.Unmarshal([]byte(messageBody), &msg); err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}

	// Validate required fields
	if msg.FCMToken == "" {
		return nil, fmt.Errorf("fcmToken is required")
	}
	if msg.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if msg.Body == "" {
		return nil, fmt.Errorf("body is required")
	}

	return &msg, nil
}
