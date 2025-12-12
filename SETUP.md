# Quick Setup Guide

## Before Deployment

1. **Create `credentials.json` for Workload Identity Federation**

   This file tells the Firebase SDK how to authenticate using your AWS Lambda role:

   ```json
   {
     "type": "external_account",
     "audience": "//iam.googleapis.com/projects/PROJECT_NUMBER/locations/global/workloadIdentityPools/aws-pool/providers/aws-provider",
     "subject_token_type": "urn:ietf:params:aws:token-type:aws4_request",
     "service_account_impersonation_url": "https://iamcredentials.googleapis.com/v1/projects/-/serviceAccounts/fcm-sender@YOUR_PROJECT_ID.iam.gserviceaccount.com:generateAccessToken",
     "token_url": "https://sts.googleapis.com/v1/token",
     "credential_source": {
       "environment_id": "aws1",
       "regional_cred_verification_url": "https://sts.{region}.amazonaws.com?Action=GetCallerIdentity&Version=2011-06-15"
     }
   }
   ```

   Replace:
   - `PROJECT_NUMBER` - Your GCP project number (not ID)
   - `YOUR_PROJECT_ID` - Your Firebase/GCP project ID

2. **Place the file in project root**
   ```bash
   # The file should be at:
   # /Users/lwinmoehein/GolandProjects/cec-notifications/credentials.json
   ```

## Deploy

```bash
# Build and deploy
make deploy

# During deployment, provide:
# - Stack name: cec-notifications
# - AWS Region: us-east-1 (or your preferred region)
# - Firebase Project ID: your-firebase-project-id
```

## Test

```bash
# Get queue URL from outputs
aws cloudformation describe-stacks \
  --stack-name cec-notifications \
  --query 'Stacks[0].Outputs[?OutputKey==`NotificationQueue`].OutputValue' \
  --output text

# Send test message
aws sqs send-message \
  --queue-url <QUEUE_URL> \
  --message-body '{"fcmToken":"YOUR_DEVICE_TOKEN","title":"Test","body":"Hello from Lambda!"}'

# View logs
make logs
```

## Troubleshooting

If you see authentication errors:
1. Verify `credentials.json` is in the deployed Lambda package
2. Check that the Lambda execution role ARN matches the one configured in GCP workload identity
3. Ensure the Firebase service account has `roles/firebase.admin` permissions
