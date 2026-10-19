package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"

	"inspirate-consulting/internal/claimtoken"
	"inspirate-consulting/internal/models"
)

const claimTokenTTL = 30 * time.Minute

var (
	sesClient        *sesv2.Client
	senderEmail      string
	appURL           string
	claimTokenSecret string
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var notification models.AccountEmailNotificationRequestBody
		if err := json.Unmarshal([]byte(record.Body), &notification); err != nil {
			log.Printf("failed to parse message %s: %v", record.MessageId, err)
			return err
		}

		token, err := claimtoken.Sign(claimTokenSecret, notification.UserEmail, notification.TempPassword, claimTokenTTL)
		if err != nil {
			log.Printf("failed to sign claim token for message %s: %v", record.MessageId, err)
			return err
		}

		claimURL := fmt.Sprintf("%s/claim-password?token=%s", appURL, url.QueryEscape(token))
		loginURL := fmt.Sprintf("%s/login", appURL)

		subject, body := formatEmail(notification, claimURL, loginURL)

		if err := sendEmail(ctx, notification.UserEmail, subject, body); err != nil {
			log.Printf("failed to send message %s: %v", record.MessageId, err)
			return err
		}

		log.Printf("sent %s email to %s", notification.EmailType, notification.UserEmail)
	}
	return nil
}

func formatEmail(n models.AccountEmailNotificationRequestBody, claimURL, loginURL string) (string, string) {
	if n.EmailType == models.PasswordReset {
		return passwordResetEmail(n.Name, claimURL, loginURL)
	}
	return accountCreatedEmail(n.Name, claimURL, loginURL)
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

	claimTokenSecret = os.Getenv("CLAIM_TOKEN_SECRET")
	if claimTokenSecret == "" {
		log.Fatal("CLAIM_TOKEN_SECRET env variable not set")
	}

	lambda.Start(Handler)
}
