package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/simpleBank/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationBearer = "Bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("Missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header ")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("Invalid authorization header")
	}

	authType := fields[0]

	if authType != authorizationBearer {
		return nil, fmt.Errorf("Unsuported uathorization type: %s", authType)
	}
	token := fields[1]
	payload, err := server.tokenMaker.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("Invalid access token")
	}
	return payload, nil
}
