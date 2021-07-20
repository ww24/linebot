package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/google/wire"
	"github.com/ww24/linebot/domain/repository"
	"golang.org/x/oauth2/google"
	f "google.golang.org/api/firestore/v1"
	"google.golang.org/api/option"
)

var Set = wire.NewSet(
	New,
	NewConversation,
	wire.Bind(new(repository.Conversation), new(*Conversation)),
)

type Client struct {
	cli *firestore.Client
}

func New(ctx context.Context) (*Client, error) {
	var opts []option.ClientOption
	var projectID string

	isEmulator := os.Getenv("FIRESTORE_EMULATOR_HOST") != ""
	if isEmulator {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
		if projectID == "" {
			projectID = "emulator"
		}
	} else {
		cred, err := google.FindDefaultCredentials(ctx, f.DatastoreScope)
		if err != nil {
			return nil, err
		}

		opts = append(opts, option.WithCredentials(cred))
		projectID = cred.ProjectID
	}

	cli, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}

	c := &Client{cli: cli}
	return c, nil
}
