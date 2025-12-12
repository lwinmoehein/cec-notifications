# CEC Notifications - Firebase FCM Lambda Service

AWS Lambda function that processes SQS events and sends Firebase Cloud Messaging (FCM) notifications using workload identity federation.

## Architecture

```
SQS Queue → Lambda Function → Firebase FCM → Mobile Devices
```

## Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/dl/)
- **AWS CLI** - [Install AWS CLI](https://aws.amazon.com/cli/)
- **AWS SAM CLI** - [Install SAM CLI](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/install-sam-cli.html)
- **Firebase Project** with FCM enabled
- **Google Cloud Workload Identity Federation** configured

## Firebase Workload Identity Federation Setup

This application uses Google Cloud Workload Identity Federation to authenticate with Firebase without storing service account keys.

### Step 1: Configure Workload Identity Pool

```bash
# Create workload identity pool
gcloud iam workload-identity-pools create aws-pool \
    --location="global" \
    --display-name="AWS Pool"

# Create workload identity provider
gcloud iam workload-identity-pools providers create-aws aws-provider \
    --location="global" \
    --workload-identity-pool="aws-pool" \
    --account-id="YOUR_AWS_ACCOUNT_ID"
```

### Step 2: Create Service Account

```bash
# Create service account for FCM
gcloud iam service-accounts create fcm-sender \
    --display-name="FCM Sender Service Account"

# Grant Firebase Admin SDK permissions
gcloud projects add-iam-policy-binding YOUR_FIREBASE_PROJECT_ID \
    --member="serviceAccount:fcm-sender@YOUR_FIREBASE_PROJECT_ID.iam.gserviceaccount.com" \
    --role="roles/firebase.admin"
```

### Step 3: Bind AWS Lambda Role

```bash
# Allow AWS Lambda role to impersonate the service account
gcloud iam service-accounts add-iam-policy-binding \
    fcm-sender@YOUR_FIREBASE_PROJECT_ID.iam.gserviceaccount.com \
    --role="roles/iam.workloadIdentityUser" \
    --member="principalSet://iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/aws-pool/attribute.aws_role/arn:aws:sts::YOUR_AWS_ACCOUNT_ID:assumed-role/LAMBDA_ROLE_NAME"
```

### Step 4: Create Credential Configuration File

Create a `credentials.json` file for the workload identity federation:

```json
{
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/aws-pool/providers/aws-provider",
  "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
  "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/fcm-sender@YOUR_FIREBASE_PROJECT_ID.iam.gserviceaccount.com:generateAccessToken",
  "token_url": "https://sts.googleapis.com/v1/token",
  "credential_source": {
    "environment_id": "aws1",
    "regional_cred_verification_url": "https://sts.{region}.amazonaws.com?Action=GetCallerIdentity&Version=2011-06-15"
  }
}
```

Place this file in your project and include it in the Lambda deployment package.

## Environment Variables

- `FIREBASE_PROJECT_ID` - Your Firebase project ID
- `GOOGLE_APPLICATION_CREDENTIALS` - Path to the credentials.json file (default: `/var/task/credentials.json`)

## SQS Message Format

Messages sent to the SQS queue should have the following JSON format:

```json
{
  "fcmToken": "device-fcm-token",
  "title": "Notification Title",
  "body": "Notification body text",
  "data": {
    "key1": "value1",
    "key2": "value2"
  }
}
```

**Required Fields:**
- `fcmToken` - FCM device token
- `title` - Notification title
- `body` - Notification message body

**Optional Fields:**
- `data` - Additional data payload (key-value pairs)

## Installation

1. **Clone and initialize:**
   ```bash
   cd /Users/lwinmoehein/GolandProjects/cec-notifications
   go mod download
   ```

2. **Place credentials file:**
   - Copy your `credentials.json` to the project root
   - Update `template.yaml` if using a different path

## Build and Deploy

### Using Makefile

```bash
# Install dependencies
make deps

# Build the Lambda function
make build

# Deploy with guided setup (first time)
make deploy

# Deploy without prompts (requires samconfig.toml)
make deploy-fast
```

### Manual Deployment

```bash
# Build
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go

# Deploy
sam build
sam deploy --guided
```

During guided deployment, provide:
- Stack name (e.g., `cec-notifications`)
- AWS Region
- Firebase Project ID
- Confirm changes and deploy

## Local Testing

### Test with SAM Local

```bash
# Invoke with sample event
make local-invoke

# Or manually
sam local invoke NotificationFunction --event events/sample-sqs-event.json
```

### Run Unit Tests

```bash
make test
```

## Monitoring

### View Logs

```bash
# Tail logs
make logs

# Or using AWS CLI
sam logs -n NotificationFunction --tail
```

### CloudWatch Metrics

Monitor the following metrics:
- Lambda invocations
- Lambda errors
- SQS queue depth
- Dead letter queue messages

## Message Flow

1. **Event Source** → Sends notification request to SQS queue
2. **SQS Queue** → Triggers Lambda function with batch of messages
3. **Lambda Function** → 
   - Parses SQS messages
   - Authenticates with Firebase using workload identity
   - Sends FCM notifications
   - Returns failed message IDs for retry
4. **Failed Messages** → Moved to Dead Letter Queue after 3 retries

## Error Handling

- **Parse Errors**: Invalid JSON or missing required fields
- **FCM Errors**: Invalid tokens, quota exceeded, etc.
- **Batch Item Failures**: Failed messages are automatically retried by SQS
- **Retry Behavior**: Messages remain in queue until successfully processed or retention period expires (14 days)

## Cleanup

```bash
# Delete the CloudFormation stack
sam delete

# Clean build artifacts
make clean
```

## Project Structure

```
.
├── main.go                 # Lambda handler
├── go.mod                  # Go dependencies
├── template.yaml           # SAM template
├── Makefile               # Build and deployment commands
├── credentials.json        # Workload identity credentials (not in git)
├── config/
│   └── firebase.go        # Firebase initialization
├── handlers/
│   ├── sqs.go             # SQS message parsing
│   └── fcm.go             # FCM message sending
└── events/
    └── sample-sqs-event.json  # Test event
```

## Troubleshooting

### Lambda can't authenticate with Firebase

- Verify workload identity pool and provider are configured correctly
- Check that Lambda execution role is allowed to impersonate the service account
- Ensure `credentials.json` is included in the deployment package

### Messages not being processed

- Check SQS queue visibility timeout (should be 3x Lambda timeout)
- Verify Lambda has permissions to read from SQS
- Check CloudWatch logs for errors

### FCM messages not delivered

- Verify Firebase project ID is correct
- Check FCM token validity
- Review Firebase console for quota limits

## License

MIT
