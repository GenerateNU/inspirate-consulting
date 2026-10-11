package sqsclient

import (
	"context"
	"log"
)

func NewSQSConnection(ctx context.Context, queueURl string) *SQSClient {
	client, err := NewSQSClient(ctx, queueURl)
	if err != nil {
		log.Fatalf("Failed to create new sqs client: %v", err)
	}

	return client

}
