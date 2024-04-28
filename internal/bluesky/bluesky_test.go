package bluesky

import (
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate_Success(t *testing.T) {
	gock.DisableNetworking()
	defer gock.Off()

	gock.New("https://bsky.social").
		Post("/xrpc/com.atproto.server.createSession").
		MatchHeader("Content-Type", "application/json").
		Reply(200).
		JSON(map[string]string{
			"Did":        "sample-did",
			"AccessJwt":  "sample-access-jwt",
			"RefreshJwt": "sample-refresh-jwt",
		})

	client, err := Authenticate("handle123", "passwordABC")

	assert.NoError(t, err)

	assert.NotNil(t, client)
	assert.Equal(t, "sample-did", client.Auth.Did)
	assert.Equal(t, "sample-access-jwt", client.Auth.AccessJwt)
	assert.Equal(t, "sample-refresh-jwt", client.Auth.RefreshJwt)

	assert.True(t, gock.IsDone(), "Not all gock interceptors were triggered")
}

func TestAuthenticate_Failure(t *testing.T) {
	gock.DisableNetworking()
	defer gock.Off()

	gock.New("https://bsky.social").
		Post("/xrpc/com.atproto.server.createSession").
		Reply(401)

	client, err := Authenticate("handle123", "passwordABC")

	assert.Error(t, err)
	assert.Nil(t, client)

	assert.True(t, gock.IsDone(), "Not all gock interceptors were triggered")
}
