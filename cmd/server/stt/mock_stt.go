package stt

import (
	"log"
	"speaking_hearts/models"
	"time"
)

// StartMockSTT simulates the Speech-to-Text processing layer.
// It listens for AudioChunks, simulates processing delay, and outputs ProcessedText.
func StartMockSTT(audioChan <-chan models.AudioChunk, textChan chan<- models.ProcessedText) {
	log.Println("Mock STT Service started")

	go func() {
		for chunk := range audioChan {
			// Simulate processing time (e.g., 500ms)
			time.Sleep(500 * time.Millisecond)

			processed := models.ProcessedText{
				OriginalChunkID: chunk.ID,
				SpeakerID:       "Mock_Speaker_1",
				OriginalLang:    chunk.LangIn,
				OriginalText:    "Hola, esta es una prueba de audio", // Simulated Spanish STT output
				Translations:    make(map[string]string),
				Timestamp:       time.Now(),
			}

			// In a real scenario, this would go to a Translator service.
			// For this simulation, we'll just send it forward.
			log.Printf("Mock STT: Processed chunk %s -> '%s'", chunk.ID, processed.OriginalText)
			textChan <- processed
		}
	}()
}
