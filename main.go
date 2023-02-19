package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type room struct {
	Label string `json:"label"`
	ID    string `json:"id"`
}

var rooms = []room{
	{Label: "Room 1", ID: "1"},
	{Label: "Room 2", ID: "2"},
	{Label: "Room 3", ID: "3"},
}

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080"},
	}))

	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoomByID)
	router.POST("/rooms", postRooms)

	router.Run("localhost:8081")
}

func getRooms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, rooms)
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

func postRooms(c *gin.Context) {
	var newRoom room

	if err := c.BindJSON(&newRoom); err != nil {
		return
	}

	rooms = append(rooms, newRoom)
	c.IndentedJSON(http.StatusCreated, newRoom)
}
