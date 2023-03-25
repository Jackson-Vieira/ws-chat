package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type connectionState struct {
	websocket *threadSafeWriter
}

type threadSafeWriter struct {
	*websocket.Conn
	sync.Mutex
}

type websocketMessage struct {
	MessageType string `json:"messageType"`
	Data        string `json:"data"`
}

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	listLock    sync.RWMutex
	connections []connectionState
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	unsafeConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	conn := &threadSafeWriter{unsafeConn, sync.Mutex{}}
	// Close the connection when the for-loop operation is finished.
	defer conn.Close() // nolint: errcheck

	listLock.Lock()
	connections = append(connections, connectionState{websocket: conn})
	listLock.Unlock()

	message := &websocketMessage{}
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		} else if err := json.Unmarshal(raw, &message); err != nil {
			log.Println(err)
			return
		}
		for _, c := range connections {
			c.websocket.WriteJSON(message)
		}
	}
}

func (t *threadSafeWriter) WriteJSON(v interface{}) error {
	t.Lock()
	defer t.Unlock()

	return t.Conn.WriteJSON(v)
}

func main() {
	http.HandleFunc("/ws", websocketHandler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
