package mock

import (
	"log"
	"speaking_hearts/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

// MockMic simulates a microphone stream with dynamic language switching.
type MockMic struct {
	CurrentLang string
	audioChan   chan<- models.AudioChunk
	mu          sync.RWMutex
}

// NewMockMic creates a new instance of MockMic.
func NewMockMic(audioChan chan<- models.AudioChunk) *MockMic {
	return &MockMic{
		CurrentLang: "es", // Default language
		audioChan:   audioChan,
	}
}

// SetLang updates the microphone's source language thread-safely.
func (m *MockMic) SetLang(lang string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	log.Printf("Mock Mic: Switching source language to [%s]", lang)
	m.CurrentLang = lang
}

// Start begins generating AudioChunks at regular intervals.
func (m *MockMic) Start() {
	ticker := time.NewTicker(3 * time.Second)
	log.Println("Mock Microphone started (Interval: 3s)")

	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			
			m.mu.RLock()
			lang := m.CurrentLang
			m.mu.RUnlock()

			chunk := models.AudioChunk{
				ID:         uuid.New().String(),
				Source:     "Mock_Mic_01",
				Timestamp:  time.Now(),
				Data:       make([]float32, 0),
				SampleRate: 44100,
				LangIn:     lang,
			}
			log.Printf("Mock Mic: Generated AudioChunk %s [Lang: %s]", chunk.ID, lang)
			m.audioChan <- chunk
		}
	}()
}
