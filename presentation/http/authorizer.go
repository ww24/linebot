package http

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
	"google.golang.org/api/idtoken"

	"github.com/ww24/linebot/internal/config"
)

type Authorizer struct {
	audience                string
	invokerServiceAccountID string
}

func NewAuthorizer(ctx context.Context, conf *config.LINEBot, cs *config.ServiceEndpoint) (*Authorizer, error) {
	serviceEndpoint, err := cs.ResolveServiceEndpoint("")
	if err != nil {
		return nil, xerrors.New("failed to get service endpoint")
	}

	return &Authorizer{
		audience:                serviceEndpoint.String(),
		invokerServiceAccountID: conf.InvokerServiceAccountID,
	}, nil
}

func (a *Authorizer) Authorize(ctx context.Context, r *http.Request) error {
	validator, err := idtoken.NewValidator(ctx)
	if err != nil {
		return xerrors.Errorf("failed to create validator to authorize: %w", err)
	}

	authorization := r.Header.Get("authorization")
	if authorization == "" {
		return xerrors.New("authorization header is empty")
	}

	authorization = strings.TrimSpace(authorization)
	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return xerrors.New("authorization header is invalid")
	}

	token := strings.TrimPrefix(authorization, prefix)
	payload, err := validator.Validate(ctx, token, a.audience)
	if err != nil {
		return xerrors.Errorf("failed to validate token: %w", err)
	}

	if payload.Issuer != "https://accounts.google.com" ||
		payload.Subject != a.invokerServiceAccountID {
		slog.WarnContext(ctx, "http: invalid token payload",
			slog.String("iss", payload.Issuer),
			slog.String("subject", payload.Subject),
		)
		return xerrors.New("token is not issued by invoker service account")
	}

	return nil
}
