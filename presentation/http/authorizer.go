package http

import (
	"context"
	"net/http"
	"strings"

	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"google.golang.org/api/idtoken"

	"github.com/ww24/linebot/domain/repository"
	"github.com/ww24/linebot/logger"
)

type Authorizer struct {
	audience                string
	invokerServiceAccountID string
}

func NewAuthorizer(ctx context.Context, conf repository.Config) (*Authorizer, error) {
	serviceEndpoint, err := conf.ServiceEndpoint("")
	if err != nil {
		return nil, xerrors.New("failed to get service endpoint")
	}

	return &Authorizer{
		audience:                serviceEndpoint.String(),
		invokerServiceAccountID: conf.InvokerServiceAccountID(),
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
		dl := logger.DefaultLogger(ctx)
		dl.Warn("invalid token payload",
			zap.String("iss", payload.Issuer),
			zap.String("subject", payload.Subject),
		)
		return xerrors.New("token is not issued by invoker service account")
	}

	return nil
}
