package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type room struct {
	Label     string `json:"label"`
	ID        string `json:"id"`
	WidgetUrl string `json:"widgetUrl"`
}

type roomLabel struct {
	Label string `json:"label"`
	ID    string `json:"id"`
}

// Global variables
var rooms = []room{}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://listening-rooms-client.onrender.com", "http://localhost:8080"},
		AllowHeaders: []string{"Origin", "Content-Type"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT"},
	}))

	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoomByID)
	router.POST("/rooms", postRoom)
	router.DELETE("/rooms/:id", deleteRoomByID)
	router.PUT("/rooms/:id", updateRoomByID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func getRooms(c *gin.Context) {
	roomLabels := []roomLabel{}

	for i := 0; i < len(rooms); i++ {
		roomLabels = append(roomLabels, roomLabel{
			Label: rooms[i].Label,
			ID:    rooms[i].ID,
		})
	}

	c.IndentedJSON(http.StatusOK, roomLabels)
}

func getRoomByID(c *gin.Context) {
	id := c.Param("id")

	for _, a := range rooms {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}

func postRoom(c *gin.Context) {
	var newRoom room

	if err := c.BindJSON(&newRoom); err != nil {
		return
	}

	newRoom.ID = uuid.New().String()

	rooms = append(rooms, newRoom)
	c.IndentedJSON(http.StatusCreated, newRoom)
}

func deleteRoomByID(c *gin.Context) {
	id := c.Param("id")

	for index, a := range rooms {
		if a.ID == id {
			rooms = append(rooms[:index], rooms[index+1:]...)
			c.IndentedJSON(http.StatusOK, gin.H{"message": "room deleted"})
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}

func updateRoomByID(c *gin.Context) {
	id := c.Param("id")

	var updatedRoom room

	if err := c.BindJSON(&updatedRoom); err != nil {
		return
	}

	for index, a := range rooms {
		if a.ID == id {
			rooms[index] = updatedRoom
			c.IndentedJSON(http.StatusOK, updatedRoom)
			return
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "room not found"})
}
