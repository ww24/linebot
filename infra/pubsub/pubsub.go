package pubsub

import (
	"context"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/google/wire"
	"golang.org/x/xerrors"

	"github.com/ww24/linebot/internal/gcp"
)

// Set provides a wire set.
var Set = wire.NewSet(
	New,
)

func New(ctx context.Context) (*pubsub.Client, error) {
	var projectID string
	isEmulator := os.Getenv("PUBSUB_EMULATOR_HOST") != ""
	if isEmulator {
		projectID = "emulator"
	} else {
		var err error
		projectID, err = gcp.ProjectID()
		if err != nil {
			return nil, xerrors.Errorf("gcp.ProjectID: %w", err)
		}
	}

	cli, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, xerrors.Errorf("pubsub.NewClient: %w", err)
	}

	return cli, nil
}
