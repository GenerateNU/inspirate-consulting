package jobs

import (
	"context"
	"inspirate-consulting/internal/models"
	"log/slog"
	"time"
)

// stub until the SQS integration is built
func SendToSqs(task models.TaskNotification, notificationType models.NotificationType) error {
	return nil
}

func (j *JobScheduler) TaskNotificationsJob() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("task notifications job panicked", "panic", r)
		}
	}()

	ctx := context.Background()

	pastDue, err := j.repo.NotificationLog.GetPastDueTasks(ctx)
	if err != nil {
		slog.Error("failed to fetch past due tasks", "error", err)
	} else {
		j.notifyTasks(ctx, pastDue, models.NotificationPastDue)
	}

	withinDue, err := j.repo.NotificationLog.GetWithinDueTasks(ctx)
	if err != nil {
		slog.Error("failed to fetch within due tasks", "error", err)
	} else {
		j.notifyTasks(ctx, withinDue, models.NotificationUpcoming)
	}
}

func (j *JobScheduler) notifyTasks(ctx context.Context, tasks []models.TaskNotification, notificationType models.NotificationType) {
	for _, task := range tasks {
		if err := SendToSqs(task, notificationType); err != nil {
			slog.Error("failed to send task notification to SQS", "task_id", task.TaskID, "error", err)
			continue
		}

		input, err := j.repo.NotificationLog.CreateNotificationLog(ctx, &models.CreateNotificationLogInput{
			UserID:           task.UserID,
			TaskID:           task.TaskID,
			NotificationType: notificationType,
			SentAt:           time.Now(),
		})
		if err != nil {
			slog.Error("failed to create notification log", "task_id", task.TaskID, "error", err)
		}

		err = SendToSqs(task, input.NotificationType)
		if err != nil {
			slog.Error("failed to send to sqs")
		}
	}
}
