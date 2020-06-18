package auth

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type Authentication struct {
	User     string
	Password string
}

func (a *Authentication) GetRequestMetadata(_ context.Context, uri ...string) (map[string]string, error) {
	log.Printf("GetRequestMetadata uri: %v\n", uri)

	return map[string]string{
		"user":     a.User,
		"password": a.Password,
	}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return true
}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("missing credentials")
	}

	var (
		appID  string
		appKey string
	)

	if val, ok := md["user"]; ok {
		appID = val[0]
	}
	if val, ok := md["password"]; ok {
		appKey = val[0]
	}

	if appID != a.User || appKey != a.Password {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}
	return nil
}
