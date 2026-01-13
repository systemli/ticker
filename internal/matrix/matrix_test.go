package matrix

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/storage"
)

func TestCreateRoom(t *testing.T) {
	// Mock server to simulate Matrix API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_matrix/client/v3/createRoom" {
			// Verify request method and headers
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

			// Send response
			response := map[string]string{
				"room_id": "!test-room-id:example.com",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else if r.URL.Path == "/_matrix/client/v3/rooms/!test-room-id:example.com/state/m.room.canonical_alias/" {
			// Return canonical alias (mautrix uses trailing slash)
			response := map[string]string{
				"alias": "#test-room:example.com",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
		} else {
			// Return empty response for other requests
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		}
	}))
	defer server.Close()

	// Create config
	cfg := config.Config{
		Matrix: config.Matrix{
			ApiUrl: server.URL,
			Token:  "test-token",
		},
	}

	// Create a ticker
	ticker := &storage.Ticker{
		Title: "Test Room",
	}

	// Test CreateRoom
	roomID, roomName, err := CreateRoom(cfg, ticker)
	assert.NoError(t, err)
	assert.Equal(t, "!test-room-id:example.com", roomID)
	assert.NotEmpty(t, roomName) // Just verify we got something back
}

func TestCreateRoom_NotConfigured(t *testing.T) {
	// Create config without Matrix configuration
	cfg := config.Config{}

	ticker := &storage.Ticker{
		Title: "Test Room",
	}

	// Test CreateRoom with unconfigured bridge
	_, _, err := CreateRoom(cfg, ticker)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "matrix bridge is not configured")
}

func TestCreateRoom_APIError(t *testing.T) {
	// Mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
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

	ticker := &storage.Ticker{
		Title: "Test Room",
	}

	// Test CreateRoom with invalid token
	_, _, err := CreateRoom(cfg, ticker)
	assert.Error(t, err)
}
