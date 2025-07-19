package bluesky

import (
	"context"
	"net/http"

	comatproto "github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/xrpc"
	"github.com/systemli/ticker/internal/logger"
)

var log = logger.GetWithPackage("bluesky")

func Authenticate(handle, password string) (*xrpc.Client, error) {
	client := &xrpc.Client{
		Client: &http.Client{},
		Host:   "https://bsky.social",
		Auth:   &xrpc.AuthInfo{Handle: handle},
	}

	auth, err := comatproto.ServerCreateSession(context.TODO(), client, &comatproto.ServerCreateSession_Input{
		Identifier: handle,
		Password:   password,
	})
	if err != nil {
		log.WithError(err).Error("failed to create session")
		return nil, err
	}

	client.Auth.Did = auth.Did
	client.Auth.AccessJwt = auth.AccessJwt
	client.Auth.RefreshJwt = auth.RefreshJwt

	return client, nil
}
