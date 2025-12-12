package config

import (
	"context"
	"fmt"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	firebaseApp     *firebase.App
	firebaseClient  *messaging.Client
	firebaseOnce    sync.Once
	firebaseInitErr error
)

// InitializeFirebase initializes the Firebase app using workload identity federation
func InitializeFirebase(ctx context.Context) (*firebase.App, error) {
	firebaseOnce.Do(func() {
		projectID := os.Getenv("FIREBASE_PROJECT_ID")
		if projectID == "" {
			firebaseInitErr = fmt.Errorf("FIREBASE_PROJECT_ID environment variable not set")
			return
		}

		// Use workload identity federation with Google Cloud credentials
		// The Lambda execution role should have permissions to assume the workload identity
		config := &firebase.Config{
			ProjectID: projectID,
		}

		// Initialize Firebase with Google Default Credentials
		// This will use the workload identity federation configured in GCP
		firebaseApp, firebaseInitErr = firebase.NewApp(ctx, config, option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")))
		if firebaseInitErr != nil {
			return
		}
	})

	return firebaseApp, firebaseInitErr
}

// GetMessagingClient returns the Firebase Messaging client
func GetMessagingClient(ctx context.Context) (*messaging.Client, error) {
	if firebaseClient != nil {
		return firebaseClient, nil
	}

	app, err := InitializeFirebase(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get messaging client: %w", err)
	}

	firebaseClient = client
	return firebaseClient, nil
}
