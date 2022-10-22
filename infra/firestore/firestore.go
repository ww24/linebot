package firestore

import (
	"context"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/wire"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/gcp"
)

// Set provides a wire set.
var Set = wire.NewSet(
	New,
	NewConversation,
	wire.Bind(new(repository.Conversation), new(*Conversation)),
	NewShopping,
	wire.Bind(new(repository.Shopping), new(*Shopping)),
	NewReminder,
	wire.Bind(new(repository.Reminder), new(*Reminder)),
)

var tracer = otel.Tracer("github.com/ww24/linebot/infra/firestore")

type Client struct {
	cli *firestore.Client
	now func() time.Time
}

func New(ctx context.Context) (*Client, error) {
	var projectID string
	isEmulator := os.Getenv("FIRESTORE_EMULATOR_HOST") != ""
	if isEmulator {
		projectID = "emulator"
	} else {
		var err error
		projectID, err = gcp.ProjectID(ctx)
		if err != nil {
			return nil, xerrors.Errorf("gcp.ProjectID: %w", err)
		}
	}

	cli, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize firestore client: %w", err)
	}

	c := &Client{
		cli: cli,
		now: time.Now,
	}
	return c, nil
}
