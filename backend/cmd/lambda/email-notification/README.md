# Email Notification Lambda

Consumes messages from the email notification SQS queue and formats a
notification email for each one. Sending through SES is not wired up yet.

## Message format

Each SQS message body is JSON matching `models.QueueEmailNotificationRequestBody`:

```json
{
  "user_email": "student@example.com",
  "task_name": "Submit FAFSA",
  "due_date": "2026-10-10T00:00:00Z",
  "is_overdue": false
}
```

## Build and deploy

Lambda runs on Amazon Linux, so the binary must be cross-compiled and
named `bootstrap`. Match `GOARCH` to the function's architecture
(`arm64`, or `amd64` for `x86_64`). From `backend/`:

```bash
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/lambda/email-notification
zip email-notification.zip bootstrap
```

Upload `email-notification.zip` in the Lambda console (Code → Upload from →
.zip file). Runtime: Amazon Linux 2023. Handler: `bootstrap`.

## Testing

Use the Lambda console's Test tab with an SQS event. The `body` field must be
the notification JSON as an escaped string.