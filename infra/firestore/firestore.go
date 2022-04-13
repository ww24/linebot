package firestore

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/google/wire"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"golang.org/x/xerrors"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/internal/gcp"
)

// Set provides a wire set.
var Set = wire.NewSet(
	New,
	NewConversation,
	wire.Bind(new(repository.Conversation), new(*Conversation)),
	NewReminder,
	wire.Bind(new(repository.Reminder), new(*Reminder)),
)

var tracer = otel.Tracer("github.com/ww24/linebot/infra/firestore")

type Client struct {
	cli *firestore.Client
}

func New(ctx context.Context) (*Client, error) {
	projectID, err := gcp.ProjectID(ctx)
	if err != nil {
		return nil, xerrors.Errorf("gcp.ProjectID: %w", err)
	}

	isEmulator := os.Getenv("FIRESTORE_EMULATOR_HOST") != ""
	if isEmulator {
		if projectID == "" {
			projectID = "emulator"
		}
	}

	opts := []option.ClientOption{
		option.WithGRPCDialOption(
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		),
		option.WithGRPCDialOption(
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
		),
	}

	cli, err := firestore.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize firestore client: %w", err)
	}

	c := &Client{cli: cli}
	return c, nil
}
