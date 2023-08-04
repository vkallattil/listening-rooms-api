package main

import "nhooyr.io/websocket"

type room struct {
	Label       string                 `json:"label"`
	ID          string                 `json:"id"`
	WidgetUrl   string                 `json:"widgetUrl"`
	SocketUsers map[string]*socketUser `json:"socketUsers"`
}

type roomLabel struct {
	Label string `json:"label"`
	ID    string `json:"id"`
}

type socketUser struct {
	ID            string
	Conn          *websocket.Conn
	CurrentRoomID string
}

type SocketMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type StringPayloadMessage struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type IntegerPayloadMessage struct {
	Type    string `json:"type"`
	Payload int    `json:"payload"`
}

type ChatMessage struct {
	SenderID   string `json:"senderID"`
	SenderName string `json:"senderName"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

type ChatPayloadMessage struct {
	Type    string        `json:"type"`
	Payload []ChatMessage `json:"payload"`
}
