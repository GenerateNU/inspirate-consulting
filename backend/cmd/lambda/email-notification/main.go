package main

import (
	"context"
	"fmt"
	"os"

	"encoding/json"
	"log"

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
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var notification models.QueueEmailNotificationRequestBody
		if err := json.Unmarshal([]byte(record.Body), &notification); err != nil {
			log.Printf("failed to parse message %s: %v", record.MessageId, err)
			return err
		}
		subject, body := formatEmail(notification)

		if err := sendEmail(ctx, notification.UserEmail, subject, body); err != nil {
			log.Printf("failed to send message %s: %v", record.MessageId, err)
			return err
		}

		log.Printf("sent email to %s for task %s", notification.UserEmail, notification.TaskName)

	}
	return nil
}

func formatEmail(n models.QueueEmailNotificationRequestBody) (string, string) {
	dueDate := n.DueDate.Format("January 2, 2006")

	if n.IsOverdue {
		subject := fmt.Sprintf("Inspirate Consulting Task Overdue: %s", n.TaskName)
		body := fmt.Sprintf("Your Insipirate Consulting task, %s, was due on %s", n.TaskName, dueDate)
		return subject, body
	}

	subject := fmt.Sprintf("Inspirate Consulting Upcoming Deadline: %s", n.TaskName)
	body := fmt.Sprintf(
		"Hi,\n\nThis is a reminder that your Inspirate Consulting task \"%s\" is due on %s.",
		n.TaskName, dueDate,
	)
	return subject, body
}

func sendEmail(ctx context.Context, to, subject, body string) error {
	_, err := sesClient.SendEmail(ctx, &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(senderEmail),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Text: &types.Content{Data: aws.String(body)},
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
	lambda.Start(Handler)
}
