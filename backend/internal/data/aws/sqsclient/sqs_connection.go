package sqsclient

import (
	"context"
	"inspirate-consulting/internal/config"
	"log"
)

func NewSQSConnection(ctx context.Context, sqsConfig config.SQS) *SQSClient {
	client, err := NewSQSClient(ctx, sqsConfig.QueueURL)
	if err != nil {
		log.Fatalf("Failed to create new sqs client: %v", err)
	}

	return client

}

func NewAccountSQSConnection(ctx context.Context, sqsConfig config.SQS) *SQSClient {
	client, err := NewSQSClient(ctx, sqsConfig.AccountQueueURL)
	if err != nil {
		log.Fatalf("Failed to create new account sqs client: %v", err)
	}

	return client
}
