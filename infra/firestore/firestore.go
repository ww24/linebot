package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/google/wire"
	"github.com/ww24/linebot/domain/repository"
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

type ClientConfig struct {
	ProjectID string
}

func New(ctx context.Context, cfg ClientConfig) (*Client, error) {
	sa := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return nil, err
	}

	cli, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	c := &Client{cli: cli}
	return c, nil
}
