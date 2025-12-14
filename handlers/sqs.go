package handlers

import (
	"encoding/json"
	"fmt"
)

// NotificationMessage represents the expected SQS message format
type NotificationMessage struct {
	ActionType string `json:"actionType"`
	FCMToken   string `json:"fcmToken"`
	TopicName  string `json:"topicName,omitempty"`

	// Now Pointers: These will be nil if the fields are omitted in the JSON payload.
	// This is ideal for SUBSCRIBE_TO_TOPIC actions.
	Title string `json:"title,omitempty"`
	Body  string `json:"body,omitempty"`

	Data map[string]string `json:"data,omitempty"`
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
	if msg.ActionType == "" {
		return nil, fmt.Errorf("action type is required")
	}

	return &msg, nil
}
