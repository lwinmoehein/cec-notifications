package handlers

import (
	"testing"
)

func TestParseSQSMessage_Valid(t *testing.T) {
	messageBody := `{
		"actionType": "SEND_SINGLE_NOTIFICATION",
		"fcmToken": "test-token-123",
		"title": "Test Title",
		"body": "Test Body",
		"data": {
			"key1": "value1",
			"key2": "value2"
		}
	}`

	msg, err := ParseSQSMessage(messageBody)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if msg.FCMToken != "test-token-123" {
		t.Errorf("Expected fcmToken 'test-token-123', got: %s", msg.FCMToken)
	}

	if msg.Title != "Test Title" {
		t.Errorf("Expected title 'Test Title', got: %s", msg.Title)
	}

	if msg.Body != "Test Body" {
		t.Errorf("Expected body 'Test Body', got: %s", msg.Body)
	}

	if len(msg.Data) != 2 {
		t.Errorf("Expected 2 data items, got: %d", len(msg.Data))
	}
}

func TestParseSQSMessage_MissingToken(t *testing.T) {
	messageBody := `{
		"actionType": "SEND_SINGLE_NOTIFICATION",
		"title": "Test Title",
		"body": "Test Body"
	}`

	_, err := ParseSQSMessage(messageBody)
	if err == nil {
		t.Error("Expected error for missing fcmToken, got nil")
	}
}

func TestParseSQSMessage_MissingTitle(t *testing.T) {
	messageBody := `{
		"actionType": "SEND_SINGLE_NOTIFICATION",
		"fcmToken": "test-token-123",
		"body": "Test Body"
	}`

	// This should pass now since title is optional (omitempty)
	msg, err := ParseSQSMessage(messageBody)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if msg.Title != "" {
		t.Errorf("Expected empty title, got: %s", msg.Title)
	}
}

func TestParseSQSMessage_MissingBody(t *testing.T) {
	messageBody := `{
		"actionType": "SEND_SINGLE_NOTIFICATION",
		"fcmToken": "test-token-123",
		"title": "Test Title"
	}`

	// This should pass now since body is optional (omitempty)
	msg, err := ParseSQSMessage(messageBody)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	if msg.Body != "" {
		t.Errorf("Expected empty body, got: %s", msg.Body)
	}
}

func TestParseSQSMessage_InvalidJSON(t *testing.T) {
	messageBody := `invalid json{`

	_, err := ParseSQSMessage(messageBody)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestParseSQSMessage_NoData(t *testing.T) {
	messageBody := `{
		"actionType": "SUBSCRIBE_TO_TOPIC",
		"fcmToken": "test-token-123",
		"topicName": "test-topic"
	}`

	msg, err := ParseSQSMessage(messageBody)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if msg.Data != nil {
		t.Errorf("Expected nil data, got: %v", msg.Data)
	}

	if msg.TopicName != "test-topic" {
		t.Errorf("Expected topicName 'test-topic', got: %s", msg.TopicName)
	}
}
