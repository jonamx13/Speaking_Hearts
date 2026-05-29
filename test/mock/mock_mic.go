package mock

import (
	"log"
	"speaking_hearts/models"
	"time"

	"github.com/google/uuid"
)

// StartMockMic simulates a microphone stream by generating AudioChunks at regular intervals.
// It sends these chunks into the provided channel, mimicking the acquisition layer.
func StartMockMic(audioChan chan<- models.AudioChunk) {
	ticker := time.NewTicker(3 * time.Second)
	log.Println("Mock Microphone started (Interval: 3s)")

	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			chunk := models.AudioChunk{
				ID:         uuid.New().String(),
				Source:     "Mock_Mic_01",
				Timestamp:  time.Now(),
				Data:       make([]float32, 0), // Empty data for simulation
				SampleRate: 44100,
				LangIn:     "es",
			}
			log.Printf("Mock Mic: Generated AudioChunk %s", chunk.ID)
			audioChan <- chunk
		}
	}()
}
