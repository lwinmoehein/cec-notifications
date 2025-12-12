package handlers

import (
	"testing"
)

func TestParseSQSMessage_Valid(t *testing.T) {
	messageBody := `{
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
		"fcmToken": "test-token-123",
		"body": "Test Body"
	}`

	_, err := ParseSQSMessage(messageBody)
	if err == nil {
		t.Error("Expected error for missing title, got nil")
	}
}

func TestParseSQSMessage_MissingBody(t *testing.T) {
	messageBody := `{
		"fcmToken": "test-token-123",
		"title": "Test Title"
	}`

	_, err := ParseSQSMessage(messageBody)
	if err == nil {
		t.Error("Expected error for missing body, got nil")
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
		"fcmToken": "test-token-123",
		"title": "Test Title",
		"body": "Test Body"
	}`

	msg, err := ParseSQSMessage(messageBody)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if msg.Data != nil {
		t.Errorf("Expected nil data, got: %v", msg.Data)
	}
}
