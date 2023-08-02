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
		OriginPatterns: []string{"https://listening-rooms-client.onrender.com", "localhost:8080"},
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

	// Defer closing the connection until the handler function returns.
	defer func() {
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

		var incomingMessage Message
		if err := json.Unmarshal(message, &incomingMessage); err != nil {
			log.Println("Error receiving message: ", err)
			return
		} else {
			log.Printf("Received message: %s\n", incomingMessage)
		}

		// if incomingMessage.Type == "SEEK" {
		// 	handleSeek(incomingMessage)
		// }

		if incomingMessage.Type == "CHANGE_ROOM" {
			handleRoomChange(&thisSocketUser, incomingMessage)
		}
	}
}

// func handleSeek(incomingMessage Message) {
// 	seekMessage := NumberPayloadMessage{
// 		Type:    "SEEK",
// 		Payload: incomingMessage.Payload.(int),
// 	}
// }

func handleRoomChange(thisSocketUser *socketUser, incomingMessage Message) {
	roomChangeMessage := StringPayloadMessage{
		Type:    "CHANGE_ROOM",
		Payload: fmt.Sprintf("%v", incomingMessage.Payload),
	}

	roomID := roomChangeMessage.Payload

	if roomID == "" {
		log.Println("room id is empty")
		if _, ok := rooms[thisSocketUser.CurrentRoomID]; ok {
			log.Println("room wasn't deleted, so no need to delete user from its current room")
			if thisSocketUser.CurrentRoomID != "" {
				delete(rooms[thisSocketUser.CurrentRoomID].SocketUsers, thisSocketUser.ID)
			}
		}
		thisSocketUser.CurrentRoomID = roomChangeMessage.Payload
		return
	}

	rooms[roomID].SocketUsers[thisSocketUser.ID] = thisSocketUser
	thisSocketUser.CurrentRoomID = roomChangeMessage.Payload

	defer func() {
		log.Println("Room change: user id:", thisSocketUser.ID, " to room id:", roomID)
	}()
}
