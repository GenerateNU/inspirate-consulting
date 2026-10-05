package jobs

import (
	"context"
	"log/slog"
	"time"
)

func (j *JobScheduler) TaskNotificationsJob() {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("task notifications job panicked")
		}
	}()

	ctx := context.Background()
	now := time.Now()

}
