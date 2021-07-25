//+build integration

package firestore

import (
	"context"
	"errors"
	"os"
	"testing"

	"google.golang.org/api/iterator"
)

const cleanupBatchSize = 100

var testClient *Client

func TestMain(m *testing.M) {
	if os.Getenv("FIRESTORE_EMULATOR_HOST") == "" {
		panic("FIRESTORE_EMULATOR_HOST must be specified in test")
	}

	var err error
	testClient, err = New(context.Background())
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func cleanupCollection(ctx context.Context, path string) error {
	colRef := testClient.cli.Collection(path)

	for {
		iter := colRef.Limit(cleanupBatchSize).Documents(ctx)
		cnt := 0

		batch := testClient.cli.Batch()
		for {
			doc, err := iter.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			cnt++
		}

		if cnt == 0 {
			return nil
		}

		if _, err := batch.Commit(ctx); err != nil {
			return err
		}
	}
}
