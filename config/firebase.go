package config

import (
	"context"
	"fmt"
	"log"
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
			log.Printf("ERROR: FIREBASE_PROJECT_ID not set")
			return
		}
		log.Printf("Firebase Project ID: %s", projectID)

		credsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		if credsPath == "" {
			firebaseInitErr = fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable not set")
			log.Printf("ERROR: GOOGLE_APPLICATION_CREDENTIALS not set")
			return
		}
		log.Printf("Credentials path: %s", credsPath)

		// Check if credentials file exists
		if _, err := os.Stat(credsPath); err != nil {
			firebaseInitErr = fmt.Errorf("credentials file not found at %s: %w", credsPath, err)
			log.Printf("ERROR: Credentials file not found at %s: %v", credsPath, err)
			return
		}
		log.Printf("✓ Credentials file exists at %s", credsPath)

		config := &firebase.Config{
			ProjectID: projectID,
		}

		// Initialize Firebase with workload identity federation credentials
		log.Printf("Initializing Firebase with workload identity federation...")
		firebaseApp, firebaseInitErr = firebase.NewApp(ctx, config, option.WithCredentialsFile(credsPath))
		if firebaseInitErr != nil {
			log.Printf("ERROR: Failed to initialize Firebase: %v", firebaseInitErr)
			return
		}
		log.Printf("✓ Firebase app initialized successfully")
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

	log.Printf("Getting Firebase Messaging client...")
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to get messaging client: %v", err)
		return nil, fmt.Errorf("failed to get messaging client: %w", err)
	}

	log.Printf("✓ Firebase Messaging client obtained successfully")
	firebaseClient = client
	return firebaseClient, nil
}
