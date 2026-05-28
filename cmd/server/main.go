package main

import (
	"log"
	"net/http"
	"speaking_hearts/cmd/server/broadcast"
	"speaking_hearts/models"
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

	// Simulated Text Ticker
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		
		counter := 0
		for {
			<-ticker.C
			counter++
			
			msg := models.ProcessedText{
				OriginalChunkID: uuid.New().String(),
				SpeakerID:       "System_Sim",
				OriginalLang:    "en",
				OriginalText:    "This is a simulated broadcast message.",
				Translations: map[string]string{
					"es": "Este es un mensaje de difusión simulado.",
					"ru": "Это симулированное широковещательное сообщение.",
					"de": "Dies ist eine simulierte Broadcast-Nachricht.",
					"fr": "Ceci est un message de diffusion simulé.",
					"zh": "这是一个模拟的广播消息。",
					"en": "This is a simulated broadcast message.",
				},
				Timestamp: time.Now(),
			}
			
			log.Printf("Simulating broadcast #%d", counter)
			manager.Broadcast <- msg
		}
	}()

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
