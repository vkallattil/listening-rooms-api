package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"nhooyr.io/websocket"
)

func getWebsocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection.
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"listening-rooms-client.onrender.com", "localhost:8080"},
	})
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create a new user.
	thisSocketUser := socketUser{
		ID:   uuid.New().String(),
		Conn: conn,
	}
	log.Println("New connection: ", thisSocketUser.ID)

	// Send the user their ID.
	idMessageBytes, err := json.Marshal(StringPayloadMessage{
		Type:    "SOCKET_ID",
		Payload: thisSocketUser.ID,
	})
	conn.Write(r.Context(), websocket.MessageText, idMessageBytes)

	// Defer closing the connection until the handler function returns.
	defer func() {
		// Remove the user from the room.
		if thisSocketUser.CurrentRoomID != "" {
			delete(rooms[thisSocketUser.CurrentRoomID].SocketUsers, thisSocketUser.ID)
		}
		conn.Close(websocket.StatusInternalError, "Internal Error")
		log.Println("Connection closed")
	}()

	// Loop Read until the connection closes.
	for {
		_, message, err := conn.Read(r.Context())
		if err != nil {
			log.Println(err)
			return
		}

		var incomingMessage SocketMessage
		if err := json.Unmarshal(message, &incomingMessage); err != nil {
			log.Println("Error receiving message: ", err)
			return
		}

		if incomingMessage.Type == "SEEK" {
			handleSeek(&thisSocketUser, incomingMessage, r)
		}

		if incomingMessage.Type == "SKIP" {
			handleSkip(&thisSocketUser, incomingMessage, r)
		}

		if incomingMessage.Type == "CHANGE_ROOM" {
			handleRoomChange(&thisSocketUser, incomingMessage, r)
		}

		if incomingMessage.Type == "CHAT" {
			handleChat(&thisSocketUser, incomingMessage, r)
		}

		if incomingMessage.Type == "PLAYBACK" {
			handlePlayback(&thisSocketUser, incomingMessage, r)
		}
	}
}

func handleSkip(thisSocketUser *socketUser, incomingMessage SocketMessage, r *http.Request) {
	skipMessage := IntegerPayloadMessage{
		Type:    "SKIP",
		Payload: int(incomingMessage.Payload.(float64)),
	}

	log.Println("Skip message: ", skipMessage.Payload)
	skipMessageBytes, err := json.Marshal(IntegerPayloadMessage{
		Type:    "SKIP",
		Payload: skipMessage.Payload,
	})
	if err != nil {
		log.Println("Error marshalling skip message: ", err)
		return
	}

	for _, socketUser := range rooms[thisSocketUser.CurrentRoomID].SocketUsers {
		if err := socketUser.Conn.Write(r.Context(), websocket.MessageText, skipMessageBytes); err != nil {
			log.Println("Error writing skip message: ", err)
			return
		}
	}
}

func handleSeek(thisSocketUser *socketUser, incomingMessage SocketMessage, r *http.Request) {
	seekMessage := IntegerPayloadMessage{
		Type:    "SEEK",
		Payload: int(incomingMessage.Payload.(float64)),
	}

	log.Println("Seek message: ", seekMessage.Payload)
	seekMessageBytes, err := json.Marshal(IntegerPayloadMessage{
		Type:    "SEEK",
		Payload: seekMessage.Payload,
	})
	if err != nil {
		log.Println("Error marshalling seek message: ", err)
		return
	}

	for _, socketUser := range rooms[thisSocketUser.CurrentRoomID].SocketUsers {
		if err := socketUser.Conn.Write(r.Context(), websocket.MessageText, seekMessageBytes); err != nil {
			log.Println("Error writing seek message: ", err)
			return
		}
	}
}

func handlePlayback(thisSocketUser *socketUser, incomingMessage SocketMessage, r *http.Request) {
	playbackMessage := StringPayloadMessage{
		Type:    "PLAYBACK",
		Payload: fmt.Sprintf("%v", incomingMessage.Payload),
	}

	log.Println("Playback message: ", playbackMessage.Payload)
	playbackMessageBytes, err := json.Marshal(StringPayloadMessage{
		Type:    "PLAYBACK",
		Payload: playbackMessage.Payload,
	})
	if err != nil {
		log.Println("Error marshalling playback message: ", err)
		return
	}

	for _, socketUser := range rooms[thisSocketUser.CurrentRoomID].SocketUsers {
		if err := socketUser.Conn.Write(r.Context(), websocket.MessageText, playbackMessageBytes); err != nil {
			log.Println("Error writing playback message: ", err)
			return
		}
	}
}

func handleChat(thisSocketUser *socketUser, incomingMessage SocketMessage, r *http.Request) {
	chatMessage := ChatMessage{
		SenderID:   thisSocketUser.ID,
		SenderName: incomingMessage.Payload.(map[string]interface{})["senderName"].(string),
		Message:    incomingMessage.Payload.(map[string]interface{})["message"].(string),
	}

	log.Println("Chat message: ", chatMessage.Message)

	chats[thisSocketUser.CurrentRoomID] = append(chats[thisSocketUser.CurrentRoomID], chatMessage)

	chatMessageBytes, err := json.Marshal(ChatPayloadMessage{
		Type:    "CHAT",
		Payload: chats[thisSocketUser.CurrentRoomID],
	})
	if err != nil {
		log.Println("Error marshalling chat message: ", err)
		return
	}

	for _, socketUser := range rooms[thisSocketUser.CurrentRoomID].SocketUsers {
		if err := socketUser.Conn.Write(r.Context(), websocket.MessageText, chatMessageBytes); err != nil {
			log.Println("Error writing chat message: ", err)
			return
		}
	}
}

func handleRoomChange(thisSocketUser *socketUser, incomingMessage SocketMessage, r *http.Request) {
	roomChangeMessage := StringPayloadMessage{
		Type:    "CHANGE_ROOM",
		Payload: fmt.Sprintf("%v", incomingMessage.Payload),
	}

	roomID := roomChangeMessage.Payload

	if roomID == "" {
		log.Println("moved to index")
		if _, ok := rooms[thisSocketUser.CurrentRoomID]; ok {
			log.Println("room wasn't deleted, so need to delete user from its current room")
			if thisSocketUser.CurrentRoomID != "" {
				delete(rooms[thisSocketUser.CurrentRoomID].SocketUsers, thisSocketUser.ID)
			}
		}
		thisSocketUser.CurrentRoomID = roomChangeMessage.Payload
		return
	}

	if thisSocketUser.CurrentRoomID != "" {
		delete(rooms[thisSocketUser.CurrentRoomID].SocketUsers, thisSocketUser.ID)
	}

	rooms[roomID].SocketUsers[thisSocketUser.ID] = thisSocketUser
	thisSocketUser.CurrentRoomID = roomChangeMessage.Payload

	initChatBytes, err := json.Marshal(ChatPayloadMessage{
		Type:    "CHAT",
		Payload: chats[thisSocketUser.CurrentRoomID],
	})
	if err != nil {
		log.Println("Error marshalling chat message: ", err)
		return
	}

	thisSocketUser.Conn.Write(r.Context(), websocket.MessageText, initChatBytes)

	log.Printf("There are currently %d users in room %s\n", len(rooms[roomID].SocketUsers), roomID)

	defer func() {
		log.Println("Room change: user id:", thisSocketUser.ID, " to room id:", roomID)
	}()
}
