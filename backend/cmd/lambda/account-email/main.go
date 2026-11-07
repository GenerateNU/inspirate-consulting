package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"inspirate-consulting/internal/models"
)

var (
	sesClient   *sesv2.Client
	senderEmail string
	appURL      string
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var notification models.AccountEmailNotificationRequestBody
		if err := json.Unmarshal([]byte(record.Body), &notification); err != nil {
			log.Printf("failed to parse message %s: %v", record.MessageId, err)
			return err
		}

		loginURL := fmt.Sprintf("%s/login", appURL)

		subject, body := formatEmail(notification, loginURL)

		if err := sendEmail(ctx, notification.UserEmail, subject, body); err != nil {
			log.Printf("failed to send message %s: %v", record.MessageId, err)
			return err
		}

		log.Printf("sent %s email to %s", notification.EmailType, notification.UserEmail)
	}
	return nil
}

func formatEmail(n models.AccountEmailNotificationRequestBody, loginURL string) (string, string) {
	if n.EmailType == models.PasswordReset {
		return passwordResetEmail(n.Name, n.TempPassword, loginURL)
	}
	return accountCreatedEmail(n.Name, n.TempPassword, loginURL)
}

func sendEmail(ctx context.Context, to, subject, htmlBody string) error {
	_, err := sesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(senderEmail),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Html: &types.Content{Data: aws.String(htmlBody)},
				},
			},
		},
	})

	return err
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Could not load AWS config: %v", err)
	}

	sesClient = sesv2.NewFromConfig(cfg)

	senderEmail = os.Getenv("SENDER_EMAIL")
	if senderEmail == "" {
		log.Fatal("SENDER_EMAIL env variable not set")
	}

	appURL = os.Getenv("APP_URL")
	if appURL == "" {
		log.Fatal("APP_URL env variable not set")
	}

	lambda.Start(Handler)
}
