package matrix

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"

	"github.com/systemli/ticker/internal/config"
	"github.com/systemli/ticker/internal/logger"
	"github.com/systemli/ticker/internal/storage"
)

var log = logger.GetWithPackage("matrix")

// getClient creates a new mautrix client with the configured credentials
func getClient(cfg config.Config) (*mautrix.Client, error) {
	if !cfg.Matrix.Enabled() {
		return nil, fmt.Errorf("matrix bridge is not configured")
	}

	client, err := mautrix.NewClient(cfg.Matrix.ApiUrl, "", cfg.Matrix.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create matrix client: %w", err)
	}
	return client, nil
}

// CreateRoom creates a new public room in Matrix using the Synapse API
func CreateRoom(cfg config.Config, t *storage.Ticker) (string, string, error) {
	client, err := getClient(cfg)
	if err != nil {
		return "", "", err
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

		roomID, err := attemptCreateRoom(client, t.Title, roomAliasName)
		if err == nil {
			// Get the canonical alias from the Matrix API
			roomName, err := getCanonicalAlias(client, roomID)
			if err != nil {
				log.WithError(err).WithField("room_id", roomID).Warn("failed to get canonical alias, using constructed alias")
				return string(roomID), roomAliasName, nil
			}
			return string(roomID), roomName, nil
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
func attemptCreateRoom(client *mautrix.Client, title, roomAliasName string) (id.RoomID, error) {
	inviteLevel := 0
	resp, err := client.CreateRoom(context.Background(), &mautrix.ReqCreateRoom{
		Name:          title,
		RoomAliasName: roomAliasName,
		Preset:        "public_chat",
		Visibility:    "public",
		PowerLevelOverride: &event.PowerLevelsEventContent{
			InvitePtr:     &inviteLevel,
			EventsDefault: 50,
		},
		InitialState: []*event.Event{
			{
				Type: event.StateEncryption,
				Content: event.Content{
					Parsed: &event.EncryptionEventContent{
						Algorithm: id.AlgorithmMegolmV1,
					},
				},
			},
		},
	})
	if err != nil {
		return "", err
	}

	log.WithField("room_id", resp.RoomID).WithField("room_name", title).WithField("room_alias", roomAliasName).Info("Matrix room created successfully")

	return resp.RoomID, nil
}

// isRoomInUseError checks if the error is due to the room alias already being taken
func isRoomInUseError(err error) bool {
	if httpErr, ok := err.(mautrix.HTTPError); ok {
		if respErr := httpErr.RespError; respErr != nil {
			return respErr.ErrCode == "M_ROOM_IN_USE"
		}
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

// getCanonicalAlias gets the canonical alias for a Matrix room
func getCanonicalAlias(client *mautrix.Client, roomID id.RoomID) (string, error) {
	var content event.CanonicalAliasEventContent
	err := client.StateEvent(context.Background(), roomID, event.StateCanonicalAlias, "", &content)
	if err != nil {
		return "", fmt.Errorf("failed to get canonical alias: %w", err)
	}

	return content.Alias.String(), nil
}

// UpdateRoomName updates the name of a Matrix room
func UpdateRoomName(cfg config.Config, roomID, name string) error {
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	_, err = client.SendStateEvent(context.Background(), id.RoomID(roomID), event.StateRoomName, "", &event.RoomNameEventContent{
		Name: name,
	})
	if err != nil {
		return fmt.Errorf("failed to update room name: %w", err)
	}

	log.WithField("room_id", roomID).WithField("room_name", name).Info("Matrix room name updated successfully")

	return nil
}

// RemoveAllMembers removes all members from a Matrix room except the bot itself
func RemoveAllMembers(cfg config.Config, roomID string) error {
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	// First, get the current user ID to avoid kicking ourselves
	whoami, err := client.Whoami(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	myUserID := whoami.UserID

	// Get all room members
	members, err := client.JoinedMembers(context.Background(), id.RoomID(roomID))
	if err != nil {
		return fmt.Errorf("failed to get room members: %w", err)
	}

	// Kick each member except ourselves
	for userID := range members.Joined {
		if userID != myUserID {
			_, err := client.KickUser(context.Background(), id.RoomID(roomID), &mautrix.ReqKickUser{
				UserID: userID,
				Reason: "Room is being deleted",
			})
			if err != nil {
				log.WithError(err).WithField("user_id", userID).Error("failed to kick user")
				// Continue kicking other users even if one fails
			} else {
				log.WithField("user_id", userID).Info("kicked user from Matrix room")
			}
		}
	}

	return nil
}

// LeaveRoom leaves a Matrix room
func LeaveRoom(cfg config.Config, roomID string) error {
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	_, err = client.LeaveRoom(context.Background(), id.RoomID(roomID))
	if err != nil {
		return fmt.Errorf("failed to leave room: %w", err)
	}

	log.WithField("room_id", roomID).Info("left Matrix room")

	return nil
}

func SendMessage(cfg config.Config, roomID, message string) error {
	client, err := getClient(cfg)
	if err != nil {
		return err
	}

	_, err = client.SendMessageEvent(context.Background(), id.RoomID(roomID), event.EventMessage, &event.MessageEventContent{
		MsgType: event.MsgText,
		Body:    message,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	log.WithField("room_id", roomID).Info("sent message to Matrix room")

	return nil
}
