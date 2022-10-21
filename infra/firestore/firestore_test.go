package firestore

import (
	"context"
	"log"
	"testing"
	"time"
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
}

func setupTestCli(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, firestoreDialTimeout)
	defer cancel()
	cli, err := New(ctx)
	if err != nil {
		return err
	}
	if _, err := cli.cli.Collections(ctx).GetAll(); err != nil {
		return err
	}
	testCli = cli
	return nil
}
