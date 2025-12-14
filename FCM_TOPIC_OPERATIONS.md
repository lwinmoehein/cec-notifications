# FCM Topic Operations - Feature Summary

## New Features Implemented

Added support for four types of FCM operations via SQS messages:

### 1. Send Single Notification
Send a direct notification to a specific device token.

### 2. Subscribe to Topic
Subscribe a device token to an FCM topic for receiving topic-based notifications.

### 3. Send to Topic
Broadcast a notification to all devices subscribed to a topic.

### 4. Unsubscribe from Topic  
Remove a device token from a topic subscription.

## Message Format

All messages require an `actionType` field:

```json
{
  "actionType": "SEND_SINGLE_NOTIFICATION" | "SUBSCRIBE_TO_TOPIC" | "SEND_TO_TOPIC_NOTIFICATION" | "UNSUBSCRIBE_FROM_TOPIC",
  "fcmToken": "device-token",      // Required for: SEND_SINGLE, SUBSCRIBE, UNSUBSCRIBE
  "topicName": "topic-name",       // Required for: SUBSCRIBE, SEND_TO_TOPIC, UNSUBSCRIBE
  "title": "Title",                // Required for: SEND_SINGLE, SEND_TO_TOPIC
  "body": "Message body",          // Required for: SEND_SINGLE, SEND_TO_TOPIC
  "data": {"key": "value"}         // Optional for all
}
```

## Implementation Details

### Handler Functions

**`SendFCMNotification()`**: Sends notification to single device
- Requires: `fcmToken`, `title`, `body`
- Optional: `data`

**`SubscribeTokenToTopic()`**: Subscribes device to topic
- Requires: `fcmToken`, `topicName`
- Returns success/failure count

**`UnsubscribeTokenFromTopic()`**: Unsubscribes device from topic
- Requires: `fcmToken`, `topicName`
- Returns success/failure count

**`SendFCMNotificationToTopic()`**: Sends to all topic subscribers
- Requires: `topicName`, `title`, `body`
- Optional: `data`

### Error Handling

- Field validation for each action type
- Detailed error logging
- Failed messages returned to SQS for retry
- Batch failure reporting enabled

## Testing

Run unit tests:
```bash
go test -v ./handlers/
```

All 6 tests passing ✅

## Deployment

```bash
sam build
sam deploy
```

## Example Usage

### Subscribe a device to "news-updates" topic
```bash
aws sqs send-message \
  --queue-url <QUEUE_URL> \
  --message-body '{"actionType":"SUBSCRIBE_TO_TOPIC","fcmToken":"device-123","topicName":"news-updates"}'
```

### Send notification to all "news-updates" subscribers
```bash
aws sqs send-message \
  --queue-url <QUEUE_URL> \
  --message-body '{"actionType":"SEND_TO_TOPIC_NOTIFICATION","topicName":"news-updates","title":"Breaking News","body":"Important update!"}'
```

### Unsubscribe from topic
```bash
aws sqs send-message \
  --queue-url <QUEUE_URL> \
  --message-body '{"actionType":"UNSUBSCRIBE_FROM_TOPIC","fcmToken":"device-123","topicName":"news-updates"}'
```

## CloudWatch Logs

Monitor operations:
```bash
aws logs tail /aws/lambda/cec-notification-processor --follow
```

Look for:
- ✓ `Subscribing token to topic: <topic>`
- ✓ `Successfully subscribed token to topic`
- ✓ `Sending notification to topic: <topic>`
- ✓ `Successfully sent FCM message to topic`
- ✓ `Unsubscribing token from topic: <topic>`
- ✓ `Successfully unsubscribed token from topic`

## Example Messages

Send To Topic:
```
{
  "actionType": "SEND_TO_TOPIC_NOTIFICATION" ,
  "topicName": "general",
  "title": "Quagmire",
  "body": "Giggity"
}
```

Send To Topic:
```
{
  "actionType": "SEND_TO_TOPIC_NOTIFICATION" ,
  "topicName": "general",
  "title": "Quagmire",
  "body": "Giggity"
}
```

Subscribe To Topic:
```
{
  "actionType": "SUBSCRIBE_TO_TOPIC",
  "fcmToken": "cRR00dudQwiBQRcFENg1uq:APA91bHE2y4cPH2yfvHypLhxVaQev8sfFTee6SYPRDGWQ_U20TNy-eiics8xdAHO_e_m_8MD62SYDGqyJSGBU_PLKNKfokvCTxNC5BKso-NRYCSuQf2_Lp4",      
  "topicName": "general"     
}
```
UnSubscribe From Topic:
```
{
  "actionType": "UNSUBSCRIBE_FROM_TOPIC",
  "fcmToken": "cRR00dudQwiBQRcFENg1uq:APA91bHE2y4cPH2yfvHypLhxVaQev8sfFTee6SYPRDGWQ_U20TNy-eiics8xdAHO_e_m_8MD62SYDGqyJSGBU_PLKNKfokvCTxNC5BKso-NRYCSuQf2_Lp4",      
  "topicName": "general"     
}
```
Send To Token:
```
{
  "actionType": "SEND_SINGLE_NOTIFICATION",
  "fcmToken": "cRR00dudQwiBQRcFENg1uq:APA91bHE2y4cPH2yfvHypLhxVaQev8sfFTee6SYPRDGWQ_U20TNy-eiics8xdAHO_e_m_8MD62SYDGqyJSGBU_PLKNKfokvCTxNC5BKso-NRYCSuQf2_Lp4",      
  "title": "Quag",
  "body": "Giggity"
}
```

