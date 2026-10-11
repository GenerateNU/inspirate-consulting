package routes

import (
	"context"
	"inspirate-consulting/internal/data/aws/sqsclient"
	"inspirate-consulting/internal/handlers/notification.go"
	"inspirate-consulting/internal/models"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func SetUpNotificationRoutes(api huma.API, sqsClient *sqsclient.SQSClient) {
	notificationHandler := notification.NewHandler(sqsClient)
	huma.Register(api, huma.Operation{
		OperationID: "queue-notification-email",
		Method:      http.MethodPost,
		Path:        "/queue/email-notification",
		Description: "Queue a single email notification to SQS",
		Tags:        []string{"Notifications"},
	}, func(ctx context.Context, input *models.QueueEmailNotificationInput) (*models.QueueEmailNotificationOutput, error) {
		messageID, err := notificationHandler.QueueEmailNotification(ctx, &input.Body)
		if(err != nil){
			return nil, err
		}

		return &models.QueueEmailNotificationOutput{Body: models.QueueEmailNotificationResponseBody{MessageID: messageID}}, nil
	})
}
