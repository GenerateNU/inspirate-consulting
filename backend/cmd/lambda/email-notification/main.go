package main

import (
	"context"
	"fmt"

	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"inspirate-consulting/internal/models"
)

func Handler(ctx context.Context, event events.SQSEvent) error {
	for _, record := range event.Records {
		var notification models.QueueEmailNotificationRequestBody
		if err := json.Unmarshal([]byte(record.Body), &notification); err != nil {
			log.Printf("failed to parse message %s: %v", record.MessageId, err)
			return err
		}
		subject, body := formatEmail(notification)
		log.Printf("formatted email for %s\nsubject: %s\nbody: %s", notification.UserEmail, subject, body)

	}
	return nil
}

func main() {
	lambda.Start(Handler)
}

func formatEmail(n models.QueueEmailNotificationRequestBody) (string, string) {
	dueDate := n.DueDate.Format("January 2, 2006")

	if n.IsOverdue {
		subject := fmt.Sprintf("Overdue: %s", n.TaskName)
		body := fmt.Sprintf("Hi,\n\nYour task \"%s\" was due on %s and is now past due.\n\nPlease complete it as soon as you can.",
			n.TaskName, dueDate)
		return subject, body
	}

	subject := fmt.Sprintf("Upcoming deadline: %s", n.TaskName)
	body := fmt.Sprintf(
		"Hi,\n\nThis is a reminder that your task \"%s\" is due on %s.",
		n.TaskName, dueDate,
	)
	return subject, body
}
