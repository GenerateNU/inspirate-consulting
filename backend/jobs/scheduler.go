package jobs

import (
	"inspirate-consulting/internal/data"
	"log"
	"log/slog"

	"github.com/robfig/cron/v3"
)

type JobScheduler struct {
	cron *cron.Cron
	repo *data.Repository
}

func NewJobScheduler(repo *data.Repository) *JobScheduler {
	return &JobScheduler{
		cron: cron.New(),
		repo: repo,
	}
}

func (j *JobScheduler) Start() {

	j.cron.Start()
	slog.Info("cron jobs have started")

	_, err := j.cron.AddFunc("0 * * * *", func() {
		log.Println("Running Task Notifications Job...")
		j.TaskNotificationsJob()
	})

	if err != nil {
		slog.Error("task notifications job failed")
	}
}

func (j *JobScheduler) Stop() {
	j.cron.Stop()
}
