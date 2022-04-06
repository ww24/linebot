package gcp

import (
	"context"
	"os"

	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
)

// ProjectID tries to detect the project ID from the environment.
// It looks in the following order:
//   1. GOOGLE_CLOUD_PROJECT envvar
//   2. ADC creds.ProjectID
func ProjectID(ctx context.Context) (string, error) {
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
