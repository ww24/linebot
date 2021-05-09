package firestore

import (
	"context"

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
	cred, err := google.FindDefaultCredentials(ctx, f.DatastoreScope)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentials(cred)
	cli, err := firestore.NewClient(ctx, cred.ProjectID, opt)
	if err != nil {
		return nil, err
	}

	c := &Client{cli: cli}
	return c, nil
}
