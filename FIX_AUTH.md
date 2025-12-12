# Firebase Authentication Fix Applied

## Problem
The `credentials.json` file was not being included in the Lambda deployment package, causing the authentication error:
```
Request had invalid authentication credentials. Expected OAuth 2 access token...
```

## Root Cause
By default, SAM only packages the Go binary (`bootstrap`) and doesn't include additional files like `credentials.json`.

## Solution Applied

### 1. Updated `template.yaml`
Added `Metadata` section to tell SAM to use Makefile for building:
```yaml
NotificationFunction:
  Type: AWS::Serverless::Function
  Metadata:
    BuildMethod: makefile  # ← This tells SAM to use Makefile
```

### 2. Updated `Makefile`
Renamed and updated the build target to copy `credentials.json`:
```makefile
build-NotificationFunction:
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap main.go
  cp credentials.json bootstrap $(ARTIFACTS_DIR)/
  chmod 644 $(ARTIFACTS_DIR)/credentials.json
```

## Next Steps

Deploy the updated code:
```bash
sam build
sam deploy
```

After deployment, the Lambda will have both:
- `/var/task/bootstrap` (the Go binary)
- `/var/task/credentials.json` (the workload identity config)

## Verification
After deployment, check CloudWatch Logs to confirm authentication works:
```bash
aws logs tail /aws/lambda/cec-notification-processor --follow
```

You should see successful FCM message sends instead of authentication errors.
