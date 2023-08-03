package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var rooms = make(map[string]*room)

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

	router.GET("/socket", func(c *gin.Context) {
		getWebsocket(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}
