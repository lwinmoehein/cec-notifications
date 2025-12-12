# Quick Fix for Firebase Authentication

## Problem
Authentication error: `Request had invalid authentication credentials`

## Solution: Use Firebase Service Account

### Step 1: Get Firebase Service Account Key

1. Go to https://console.firebase.google.com/
2. Select project: **cec-app-6b091**
3. Click **⚙️ Settings** → **Project Settings** → **Service Accounts** tab
4. Click **Generate New Private Key** button
5. Download the JSON file

### Step 2: Replace credentials.json

**IMPORTANT**: Replace your current `credentials.json` file with the one you just downloaded.

The file should start with:
```json
{
  "type": "service_account",
  "project_id": "cec-app-6b091",
  ...
}
```

**NOT** `"type": "external_account"` (that's for workload identity).

### Step 3: Deploy

```bash
sam build
sam deploy
```

### Step 4: Test

Send a test message:
```bash
# Get your queue URL from CloudFormation outputs
aws cloudformation describe-stacks \
  --query 'Stacks[?StackName!=`null`]|[?Outputs!=`null`]|[0].Outputs[?OutputKey==`NotificationQueue`].OutputValue' \
  --output text

# Send test (replace QUEUE_URL and FCM_TOKEN)
aws sqs send-message \
  --queue-url <QUEUE_URL> \
  --message-body '{"fcmToken":"<YOUR_FCM_TOKEN>","title":"Test","body":"It works!"}'
```

### Step 5: Check Logs

```bash
aws logs tail /aws/lambda/cec-notification-processor --follow
```

You should see:
- ✅ `Firebase Project ID: cec-app-6b091`
- ✅ `Credentials file exists at /var/task/credentials.json`
- ✅ `Firebase app initialized successfully`
- ✅ `Successfully sent FCM message`

No more authentication errors! 🎉
