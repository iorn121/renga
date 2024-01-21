package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// 初期化
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Participant struct {
	ID       string
	Username string
	RoomID   string
}

type Room struct {
	ID              string
	Participants    []*websocket.Conn
	MaxParticipants int
	ReadPoems       []string
	Order           []int
}

func NewRoom(id string) *Room {
	return &Room{
		ID:              id,
		Participants:    make([]*websocket.Conn, 0),
		MaxParticipants: 5,
		ReadPoems:       make([]string, 0),
		Order:           make([]int, 0),
	}
}

func (r *Room) Join(conn *websocket.Conn) error {
	if len(r.Participants) >= r.MaxParticipants {
		return fmt.Errorf("room is full")
	}
	r.Participants = append(r.Participants, conn)
	return nil
}

func main() {
	room := NewRoom("room1")

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil)
		err := room.Join(conn)
		if err != nil {
			fmt.Println(err)
			return
		}

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			for _, client := range room.Participants {
				if err = client.WriteMessage(msgType, msg); err != nil {
					return
				}
			}
		}
	})

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./src/main.jsx")
	// })
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	println("Server started at http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
