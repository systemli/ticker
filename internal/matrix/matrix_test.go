package matrix

import (
"encoding/json"
"net/http"
"net/http/httptest"
"testing"

"github.com/stretchr/testify/assert"
"github.com/systemli/ticker/internal/config"
)

func TestCreateRoom(t *testing.T) {
// Mock server to simulate Synapse API
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// Verify request method and path
assert.Equal(t, http.MethodPost, r.Method)
assert.Equal(t, "/_matrix/client/v3/createRoom", r.URL.Path)

// Verify authorization header
assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

// Verify request body
var req CreateRoomRequest
err := json.NewDecoder(r.Body).Decode(&req)
assert.NoError(t, err)
assert.Equal(t, "test-room", req.Name)
assert.Equal(t, "public_chat", req.Preset)
assert.Equal(t, "public", req.Visibility)
assert.Equal(t, "test-room", req.RoomAliasName)

// Send response
response := CreateRoomResponse{
RoomID: "!test-room-id:example.com",
}
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(response)
}))
defer server.Close()

// Create config
cfg := config.Config{
Matrix: config.Matrix{
ApiUrl: server.URL,
Token:  "test-token",
},
}

// Test CreateRoom
roomID, err := CreateRoom(cfg, "test-room")
assert.NoError(t, err)
	assert.Equal(t, "!test-room-id:example.com", roomID)
}

func TestCreateRoom_NotConfigured(t *testing.T) {
// Create config without Matrix configuration
cfg := config.Config{}

// Test CreateRoom with unconfigured bridge
_, err := CreateRoom(cfg, "test-room")
assert.Error(t, err)
assert.Contains(t, err.Error(), "matrix bridge is not configured")
}

func TestCreateRoom_APIError(t *testing.T) {
// Mock server that returns an error
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
w.WriteHeader(http.StatusUnauthorized)
w.Write([]byte(`{"errcode":"M_FORBIDDEN","error":"Invalid access token"}`))
}))
defer server.Close()

// Create config
cfg := config.Config{
Matrix: config.Matrix{
ApiUrl: server.URL,
Token:  "invalid-token",
},
}

// Test CreateRoom with invalid token
_, err := CreateRoom(cfg, "test-room")
assert.Error(t, err)
assert.Contains(t, err.Error(), "API request failed with status 401")
}
