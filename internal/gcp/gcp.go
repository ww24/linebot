package gcp

import (
	"context"
	"os"
	"sync"
	"time"

	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

const fetchTimeout = 5 * time.Second

//nolint:gochecknoglobals
var projectID = sync.OnceValues(withTimeout(getProjectID, fetchTimeout))

// ProjectID tries to detect the project ID from the environment.
// It looks in the following order:
//  1. GOOGLE_CLOUD_PROJECT envvar
//  2. ADC creds.ProjectID
func ProjectID() (string, error) {
	return projectID()
}

func withTimeout(f func(ctx context.Context) (string, error), d time.Duration) func() (string, error) {
	return func() (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), d)
		defer cancel()
		return f(ctx)
	}
}

func getProjectID(ctx context.Context) (string, error) {
	if projectID := os.Getenv("GOOGLE_CLOUD_PROJECT"); projectID != "" {
		return projectID, nil
	}

	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return "", xerrors.Errorf("fetching default credentials: %w", err)
	}

	if creds.ProjectID == "" {
		return "", xerrors.New("unable to detect projectID")
	}

	return creds.ProjectID, nil
}
