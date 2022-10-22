package firestore

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"go.opentelemetry.io/otel"
	"google.golang.org/api/iterator"

	"github.com/ww24/linebot/domain/model"
)

const firestoreDialTimeout = 5 * time.Second

var testCli *Client

func TestMain(m *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := setupTestCli(ctx); err != nil {
		log.Fatalf("failed to setup firestore client: %v", err)
	}

	m.Run()

	bw := testCli.cli.BulkWriter(ctx)
	defer bw.End()
	conv := NewConversation(testCli)
	ids, err := removeAllDocuments(bw, conv.conversations().DocumentRefs(ctx))
	if err != nil {
		panic(err)
	}
	for _, id := range ids {
		conversationID := model.ConversationID(id)
		s := NewShopping(conv).shopping(conversationID)
		if _, err := removeAllDocuments(bw, s.DocumentRefs(ctx)); err != nil {
			panic(err)
		}
		r := NewReminder(conv).reminder(conversationID)
		if _, err := removeAllDocuments(bw, r.DocumentRefs(ctx)); err != nil {
			panic(err)
		}
	}
}

func setupTestCli(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, firestoreDialTimeout)
	defer cancel()
	cli, err := New(ctx, otel.GetTracerProvider())
	if err != nil {
		return err
	}
	if _, err := cli.cli.Collections(ctx).GetAll(); err != nil {
		return err
	}
	testCli = cli
	return nil
}

func removeAllDocuments(bw *firestore.BulkWriter, refItr *firestore.DocumentRefIterator) ([]string, error) {
	ids := make([]string, 0)
	for {
		ref, err := refItr.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if _, err := bw.Delete(ref, firestore.Exists); err != nil {
			return nil, err
		}
		ids = append(ids, ref.ID)
	}
	return ids, nil
}

func (c *Client) clone() *Client {
	cli := *c
	return &cli
}
