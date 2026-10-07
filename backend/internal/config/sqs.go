package config

type SQS struct{
	QueueURL string `env: "SQS_URL, default=https://sqs.us-east-1.amazonaws.com/478867930449/inspirate-email-notification-queue"`
}