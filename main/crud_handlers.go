package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getRooms(c *gin.Context) {
	roomLabels := []roomLabel{}

	for _, room := range rooms {
		roomLabels = append(roomLabels, roomLabel{
			Label: room.Label,
			ID:    room.ID,
		})
	}

	c.IndentedJSON(http.StatusOK, roomLabels)
}

func getRoomByID(c *gin.Context) {
	id := c.Param("id")

	for _, room := range rooms {
		if room.ID == id {
			c.IndentedJSON(http.StatusOK, room)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}

func postRoom(c *gin.Context) {
	var newRoom room

	if err := c.BindJSON(&newRoom); err != nil {
		log.Println(err)
		return
	}

	newRoom.ID = uuid.New().String()
	newRoom.SocketUsers = make(map[string]*socketUser)

	rooms[newRoom.ID] = &newRoom
	c.IndentedJSON(http.StatusCreated, newRoom)
}

func deleteRoomByID(c *gin.Context) {
	id := c.Param("id")

	if rooms[id] != nil {
		delete(rooms, id)
		c.IndentedJSON(http.StatusOK, gin.H{"message": "room deleted"})
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}

func updateRoomByID(c *gin.Context) {
	id := c.Param("id")

	var updatedRoom room

	if err := c.BindJSON(&updatedRoom); err != nil {
		log.Println(err)
		return
	}

	if rooms[id] != nil {
		rooms[id] = &updatedRoom
		c.IndentedJSON(http.StatusOK, updatedRoom)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}
