package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
)

var log = logger.GetWithPackage("matrix")

// CreateRoomRequest represents the request payload for creating a Matrix room
type CreateRoomRequest struct {
	Name                      string                   `json:"name"`
	RoomAliasName             string                   `json:"room_alias_name"`
	Preset                    string                   `json:"preset"`
	Visibility                string                   `json:"visibility"`
	PowerLevelContentOverride map[string]interface{}   `json:"power_level_content_override,omitempty"`
	InitialState              []map[string]interface{} `json:"initial_state,omitempty"`
}

// CreateRoomResponse represents the response from creating a Matrix room
type CreateRoomResponse struct {
	RoomID string `json:"room_id"`
}

// MatrixErrorResponse represents an error response from the Matrix API
type MatrixErrorResponse struct {
	ErrCode string `json:"errcode"`
	Error   string `json:"error"`
}

// MembersResponse represents the response from getting room members
type MembersResponse struct {
	Chunk []MemberEvent `json:"chunk"`
}

// MemberEvent represents a single member event
type MemberEvent struct {
	Type     string            `json:"type"`
	StateKey string            `json:"state_key"`
	Content  map[string]string `json:"content"`
}

// KickRequest represents the request to kick a user from a room
type KickRequest struct {
	UserID string `json:"user_id"`
	Reason string `json:"reason"`
}

// CreateRoom creates a new public room in Matrix using the Synapse API
func CreateRoom(cfg config.Config, t *storage.Ticker) (string, string, error) {
	if !cfg.Matrix.Enabled() {
		return "", "", fmt.Errorf("matrix bridge is not configured")
	}

	// Sanitize room name: convert to ASCII-only, remove spaces and special characters
	baseRoomName := sanitizeRoomName(t.Title)

	// Try to create the room, incrementing a suffix if the alias is already taken
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		roomAliasName := baseRoomName
		if i > 0 {
			roomAliasName = fmt.Sprintf("%s-%d", baseRoomName, i)
		}

		roomID, err := attemptCreateRoom(cfg, t.Title, roomAliasName)
		if err == nil {
			return roomID, roomAliasName, nil
		}

		// Check if the error is due to room alias already being taken
		if !isRoomInUseError(err) {
			return "", "", err
		}

		// If it's a room in use error, try again with incremented suffix
		log.WithField("room_alias", roomAliasName).Debug("Room alias already taken, trying with incremented suffix")
	}

	return "", "", fmt.Errorf("failed to create room after %d attempts: all aliases taken", maxRetries)
}

// attemptCreateRoom attempts to create a Matrix room with the given alias name
func attemptCreateRoom(cfg config.Config, title, roomAliasName string) (string, error) {
	url := fmt.Sprintf("%s/_matrix/client/v3/createRoom", cfg.Matrix.ApiUrl)

	requestBody := CreateRoomRequest{
		Name:          title,
		RoomAliasName: roomAliasName,
		Preset:        "public_chat",
		Visibility:    "public",
		// Set default power levels: allow anyone to invite and restrict sending messages to moderators and above
		PowerLevelContentOverride: map[string]interface{}{
			"invite":         0,
			"events_default": 50,
		},
		// Enable end-to-end encryption for the room
		InitialState: []map[string]interface{}{
			{
				"type": "m.room.encryption",
				"content": map[string]interface{}{
					"algorithm": "m.megolm.v1.aes-sha2",
				},
			},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Matrix.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Try to parse the error response
		var matrixErr MatrixErrorResponse
		if err := json.Unmarshal(body, &matrixErr); err == nil {
			return "", &matrixError{
				StatusCode: resp.StatusCode,
				ErrCode:    matrixErr.ErrCode,
				Message:    matrixErr.Error,
			}
		}
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response CreateRoomResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	log.WithField("room_id", response.RoomID).WithField("room_name", title).WithField("room_alias", roomAliasName).Info("Matrix room created successfully")

	return response.RoomID, nil
}

// matrixError represents a structured Matrix API error
type matrixError struct {
	StatusCode int
	ErrCode    string
	Message    string
}

func (e *matrixError) Error() string {
	return fmt.Sprintf("Matrix API error (status %d): %s - %s", e.StatusCode, e.ErrCode, e.Message)
}

// isRoomInUseError checks if the error is due to the room alias already being taken
func isRoomInUseError(err error) bool {
	if matrixErr, ok := err.(*matrixError); ok {
		return matrixErr.ErrCode == "M_ROOM_IN_USE"
	}
	return false
}

// sanitizeRoomName converts a room name to an ASCII-only variant without spaces
func sanitizeRoomName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Remove non-ASCII characters
	asciiOnly := regexp.MustCompile(`[^a-z0-9-_]`)
	name = asciiOnly.ReplaceAllString(name, "")

	// Trim leading and trailing special characters
	name = strings.Trim(name, "-_")

	// If the name is empty after sanitization, use a default
	if name == "" {
		name = "ticker-room"
	}

	return name
}

// RemoveAllMembers removes all members from a Matrix room except the bot itself
func RemoveAllMembers(cfg config.Config, roomID string) error {
	if !cfg.Matrix.Enabled() {
		return fmt.Errorf("matrix bridge is not configured")
	}

	// First, get the current user ID to avoid kicking ourselves
	myUserID, err := getCurrentUser(cfg)
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	log.Debug("Current Matrix user ID: ", myUserID)

	// Get all room members
	members, err := getRoomMembers(cfg, roomID)
	if err != nil {
		return fmt.Errorf("failed to get room members: %w", err)
	}
	log.Debug("Current Matrix room members: ", members)

	// Kick each member except ourselves
	for _, member := range members {
		if member.StateKey != myUserID {
			err := kickUser(cfg, roomID, member.StateKey)
			if err != nil {
				log.WithError(err).WithField("user_id", member.StateKey).Error("failed to kick user")
				// Continue kicking other users even if one fails
			} else {
				log.WithField("user_id", member.StateKey).Info("kicked user from Matrix room")
			}
		}
	}

	return nil
}

// getCurrentUser gets the user ID of the bot
func getCurrentUser(cfg config.Config) (string, error) {
	url := fmt.Sprintf("%s/_matrix/client/v3/account/whoami", cfg.Matrix.ApiUrl)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Matrix.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response struct {
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.UserID, nil
}

// getRoomMembers gets all members of a Matrix room
func getRoomMembers(cfg config.Config, roomID string) ([]MemberEvent, error) {
	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/members", cfg.Matrix.ApiUrl, roomID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Matrix.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response MembersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response.Chunk, nil
}

// kickUser kicks a user from a Matrix room
func kickUser(cfg config.Config, roomID, userID string) error {
	url := fmt.Sprintf("%s/_matrix/client/r0/rooms/%s/kick", cfg.Matrix.ApiUrl, roomID)

	kickReq := KickRequest{
		UserID: userID,
		Reason: "Room is being deleted",
	}

	jsonData, err := json.Marshal(kickReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Matrix.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// LeaveRoom leaves a Matrix room
func LeaveRoom(cfg config.Config, roomID string) error {
	if !cfg.Matrix.Enabled() {
		return fmt.Errorf("matrix bridge is not configured")
	}

	url := fmt.Sprintf("%s/_matrix/client/v3/rooms/%s/leave", cfg.Matrix.ApiUrl, roomID)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfg.Matrix.Token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	log.WithField("room_id", roomID).Info("left Matrix room")

	return nil
}
