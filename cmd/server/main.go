package main

import (
	"log"
	"net/http"
	"speaking_hearts/cmd/server/broadcast"
	"speaking_hearts/cmd/server/config"
	"speaking_hearts/cmd/server/stt"
	"speaking_hearts/cmd/server/translate"
	"speaking_hearts/cmd/server/tts"
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
	textChan := make(chan models.ProcessedText, 100)
	ttsChan := make(chan models.ProcessedText, 100)

	// Start Mock Microphone (Acquisition Layer)
	mock.StartMockMic(audioChan)

	// Initialize the STT Engine (Whisper Skeleton)
	whisper := stt.NewWhisperEngine("models/whisper-base")

	// Start the STT Worker Pool (Processing Layer)
	sttPool := stt.NewWorkerPool(3, audioChan, textChan, whisper)
	sttPool.Start()

	// Initialize the Translation Service and Router (Processing Layer)
	nllb := translate.NewTranslatorService("models/nllb/distilled-600M")
	
	// Load routing rules from configuration
	rules := config.GetDefaultRoutingRules()
	router := translate.NewLanguageRouter(nllb, rules)

	// Wire the STT output to the Translation Router
	go func() {
		log.Println("Translation Router started")
		for processedText := range textChan {
			// Apply translation rules
			router.RouteProcess(&processedText)
			// Forward to TTS pipeline
			ttsChan <- processedText
		}
	}()

	// Initialize the TTS Engine and Worker Pool (Processing Layer)
	ttsEngine := tts.NewMockEngine()
	ttsPool := tts.NewWorkerPool(2, ttsChan, manager.Broadcast, ttsEngine)
	ttsPool.Start()

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
