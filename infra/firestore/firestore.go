package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/google/wire"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
	f "google.golang.org/api/firestore/v1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/ww24/linebot/domain/repository"
)

// Set provides a wire set.
//nolint: gochecknoglobals
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
			return nil, xerrors.Errorf("failed to find default credentials: %w", err)
		}

		opts = append(opts, option.WithCredentials(cred))
		projectID = cred.ProjectID
	}

	opts = append(opts,
		option.WithGRPCDialOption(
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		),
		option.WithGRPCDialOption(
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		),
	)

	cli, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize firestore client: %w", err)
	}

	c := &Client{cli: cli}
	return c, nil
}
