package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

// Type definitions
type user struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
}

type room struct {
	Label   string `json:"label"`
	ID      string `json:"id"`
	Private bool   `json:"isPrivate"`
	Users   []user `json:"users"`
	SongURL string `json:"songUrl"`
}

type roomLabel struct {
	Label     string `json:"label"`
	ID        string `json:"id"`
	Private   bool   `json:"isPrivate"`
	UserCount uint8  `json:"userCount"`
}

// Global variables

var rooms = []room{
	{Label: "Room 1", ID: "1", Private: true, Users: []user{{UserName: "johndoe", ID: "1"}}, SongURL: "https://soundcloud.com/shaggy_svare_ganj420/travis-scott-yeah-yeah-feat-young-thugbass-boosted?si=1acb85f388a64d2bb97f0b3bc971003f&utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing"},
	{Label: "Room 2", ID: "2", Private: false, Users: []user{{ID: "2", UserName: "vkallattil"}}, SongURL: "https://soundcloud.com/premierhiphopdaily/toomanychances?si=00da0a6bc9354fb0a76b0197d345f813&utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing"},
	{Label: "Room 3", ID: "3", Private: false, Users: []user{{ID: "3", UserName: "amyadams"}}, SongURL: "https://soundcloud.com/anthony-bell-463042235/21-savage-gang-shit-knife-talk-og?si=c6ae54ffd19b482fa32da7cd999999c3&utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing"},
	{Label: "Room 4", ID: "4", Private: true, Users: []user{{ID: "4", UserName: "pointerfunk"}}, SongURL: "https://soundcloud.com/user-300275914/21-savage-yea-yea-prod-pierre-bourne?si=6700eb57a94b410688efebbfc74a9779&utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing"},
}

func main() {
	router := gin.Default()
	socketServer := socketio.NewServer(nil)

	socketServer.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	socketServer.OnError("/", func(s socketio.Conn, e error) {
		log.Println("error:", e)
	})

	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	socketServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socketServer.Close()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:8080"},
	}))

	router.GET("/rooms", getRooms)
	router.GET("/rooms/:id", getRoomByID)
	router.POST("/rooms", postRoom)
	router.GET("/socket.io/*any", gin.WrapH(socketServer))
	router.POST("/socket.io/*any", gin.WrapH(socketServer))

	if err := router.Run("localhost:8081"); err != nil {
		log.Fatal("failed run app: ", err)
	}
}

func getRooms(c *gin.Context) {
	var roomLabels []roomLabel

	for i := 0; i < len(rooms); i++ {
		roomLabels = append(roomLabels, roomLabel{
			Label: rooms[i].Label,
			ID:    rooms[i].ID, Private: rooms[i].Private,
			UserCount: uint8(len(rooms[i].Users))})
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

	rooms = append(rooms, newRoom)
	c.IndentedJSON(http.StatusCreated, newRoom)
}
