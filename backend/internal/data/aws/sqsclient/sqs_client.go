package sqsclient

import (
	"context"
	"fmt"
	"log"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	
	"inspirate-consulting/internal/models"
	
)

// SQSAPI is the subset of the AWS SQS client used by SQSClient, so it can be faked in tests.
type SQSAPI interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SQSClient struct {
	sqsClient SQSAPI
	queueURL  string
}

// NewSQSClientWithAPI builds an SQSClient from an existing SQS API implementation.
func NewSQSClientWithAPI(api SQSAPI, queueURL string) *SQSClient {
	return &SQSClient{sqsClient: api, queueURL: queueURL}
}

func NewSQSClient(ctx context.Context, queueURL string) (*SQSClient, error) {
	if queueURL == "" {
		return nil, fmt.Errorf("sqsclietn: queueURL required")
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client := sqs.NewFromConfig(cfg)

	return &SQSClient{sqsClient: client, queueURL: queueURL}, nil
}

func (s *SQSClient) SendEmailNotification(ctx context.Context, notification *models.QueueEmailNotificationRequestBody) (*string, error) {
	body, err := json.Marshal(notification)
	if(err != nil){
		return nil, err
	}
	result, err := s.sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(string(body)),
		QueueUrl : aws.String(s.queueURL),
	})
	if(err != nil){
		return nil, err
	}
	return result.MessageId, nil

}

