package main

import (
	"log"
	"net/http"
	"speaking_hearts/cmd/server/broadcast"
	"speaking_hearts/cmd/server/stt"
	"speaking_hearts/models"
	"speaking_hearts/test/mock"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for local simulation
	},
}

func serveWs(manager *broadcast.Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// Extract preferred language from query params, default to "en"
	lang := r.URL.Query().Get("lang")
	if lang == "" {
		lang = "en"
	}

	client := &broadcast.Client{
		ID:            uuid.New().String(),
		Conn:          conn,
		Type:          "listener",
		PreferredLang: lang,
		LastSeen:      time.Now(),
		Send:          make(chan []byte, 256),
	}

	manager.Register <- client

	// Start the reader and writer loops in separate goroutines
	go client.WritePump()
	go client.ReadPump(manager)
}

func main() {
	manager := broadcast.NewManager()
	go manager.Run()

	// Initialize channels for the processing pipeline
	audioChan := make(chan models.AudioChunk, 100)

	// Start Mock Microphone (Acquisition Layer)
	mock.StartMockMic(audioChan)

	// Start Mock STT (Processing Layer)
	// We wire STT directly to the manager's broadcast channel for now.
	stt.StartMockSTT(audioChan, manager.Broadcast)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(manager, w, r)
	})

	fs := http.FileServer(http.Dir("./web"))
	http.Handle("/", fs)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
