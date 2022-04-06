package scheduler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/go-oss/scheduler"
	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/model"
	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/gcp"
)

// Set provides a wire set.
var Set = wire.NewSet(
	New,
	wire.Bind(new(repository.ScheduleSynchronizer), new(*Scheduler)),
)

const (
	taskPrefix       = "reminder-"
	reminderEndpoint = "/reminder"
)

type Scheduler struct {
	cli       *cloudtasks.Client
	projectID string
	location  string
	queue     string
	endpoint  *url.URL
}

func New(ctx context.Context, conf repository.Config) (*Scheduler, error) {
	projectID, err := gcp.ProjectID(ctx)
	if err != nil {
		return nil, xerrors.Errorf("gcp.ProjectID: %w", err)
	}

	cli, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize CloudTasks client: %w", err)
	}

	endpoint, err := conf.ServiceEndpoint(reminderEndpoint)
	if err != nil {
		return nil, xerrors.Errorf("failed to get service endpoint: %w", err)
	}

	return &Scheduler{
		cli:       cli,
		projectID: projectID,
		location:  conf.CloudTasksLocation(),
		queue:     conf.CloudTasksQueue(),
		endpoint:  endpoint,
	}, nil
}

func (s *Scheduler) prefix(conversationID model.ConversationID) string {
	return taskPrefix + conversationID.String() + "-"
}

func (s *Scheduler) Sync(ctx context.Context, conversationID model.ConversationID, items model.ReminderItems, t time.Time) error {
	prefix := s.prefix(conversationID)
	sc := scheduler.New(s.cli, s.projectID, s.location, s.queue, prefix)

	tasks := make([]*scheduler.Task, 0, len(items))
	for _, item := range items {
		next, err := item.Scheduler.Next(t)
		if err != nil {
			continue
		}
		data, err := json.Marshal(item.IDJSON())
		if err != nil {
			return xerrors.Errorf("failed to marshal json: %w", err)
		}
		req, err := http.NewRequest(http.MethodPost, s.endpoint.String(), bytes.NewReader(data))
		if err != nil {
			return xerrors.Errorf("failed to create http request: %w", err)
		}
		req.Header.Set("content-type", "application/json")
		tasks = append(tasks, &scheduler.Task{
			QueuePath:   scheduler.QueuePath(s.projectID, s.location, s.queue),
			Prefix:      prefix,
			ID:          string(item.ID),
			ScheduledAt: next,
			Request:     req,
			Version:     1,
		})
	}

	if err := sc.Sync(ctx, tasks); err != nil {
		return xerrors.Errorf("failed to sync tasks: %w", err)
	}

	return nil
}
