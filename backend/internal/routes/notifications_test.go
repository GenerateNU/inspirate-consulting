package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"inspirate-consulting/internal/config"
	"inspirate-consulting/internal/data"
	"inspirate-consulting/internal/data/aws/sqsclient"
	"inspirate-consulting/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSQS records the last message sent and returns a canned response.
type fakeSQS struct {
	lastInput *sqs.SendMessageInput
	messageID string
	err       error
}

func (f *fakeSQS) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	f.lastInput = params
	if f.err != nil {
		return nil, f.err
	}
	return &sqs.SendMessageOutput{MessageId: aws.String(f.messageID)}, nil
}

func setupTestAppWithNotifications(fake *fakeSQS) (*fiber.App, error) {
	cfg := config.Config{
		TestMode: true,
	}
	sqsClient := sqsclient.NewSQSClientWithAPI(fake, "test-queue-url")
	app, _, err := SetupApp(cfg, &data.Repository{}, sqsClient, nil)
	return app, err
}

func TestRoute_QueueEmailNotification(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		fake := &fakeSQS{messageID: "msg-123"}

		app, err := setupTestAppWithNotifications(fake)
		require.NoError(t, err)

		payload := map[string]any{
			"user_email": "student@example.com",
			"task_name":  "Submit essay",
			"due_date":   "2026-10-15T00:00:00Z",
			"is_overdue": false,
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/queue/email-notification", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var output models.QueueEmailNotificationResponseBody
		err = json.Unmarshal(respBody, &output)
		require.NoError(t, err)
		require.NotNil(t, output.MessageID)
		assert.Equal(t, "msg-123", *output.MessageID)

		// Verify the message was sent to the configured queue with the request body
		require.NotNil(t, fake.lastInput)
		assert.Equal(t, "test-queue-url", *fake.lastInput.QueueUrl)

		var sent models.QueueEmailNotificationRequestBody
		err = json.Unmarshal([]byte(*fake.lastInput.MessageBody), &sent)
		require.NoError(t, err)
		assert.Equal(t, "student@example.com", sent.UserEmail)
		assert.Equal(t, "Submit essay", sent.TaskName)
	})

	t.Run("sqs error -> 500", func(t *testing.T) {
		t.Parallel()

		fake := &fakeSQS{err: errors.New("sqs unavailable")}

		app, err := setupTestAppWithNotifications(fake)
		require.NoError(t, err)

		payload := map[string]any{
			"user_email": "student@example.com",
			"task_name":  "Submit essay",
			"due_date":   "2026-10-15T00:00:00Z",
			"is_overdue": true,
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/queue/email-notification", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("validation error -> missing required field", func(t *testing.T) {
		t.Parallel()

		fake := &fakeSQS{messageID: "msg-123"}
		// SQS should not be called when Huma validation fails.

		app, err := setupTestAppWithNotifications(fake)
		require.NoError(t, err)

		payload := map[string]any{
			"user_email": "student@example.com",
			// missing every other required field
		}
		bodyBytes, err := json.Marshal(payload)
		require.NoError(t, err)

		req, err := http.NewRequest(http.MethodPost, "/queue/email-notification", bytes.NewReader(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		defer func() { _ = resp.Body.Close() }()

		assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)
		assert.Nil(t, fake.lastInput)
	})
}
