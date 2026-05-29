package stt

import (
	"log"
	"speaking_hearts/models"
	"time"
)

// STTEngine defines the interface for speech-to-text processing.
// This allows the WorkerPool to remain agnostic of the specific AI implementation.
type STTEngine interface {
	Transcribe(chunk models.AudioChunk) (string, error)
}

// WhisperEngine is the wrapper for the faster-whisper STT engine.
// Currently, it acts as a simulation placeholder for the real CGo bindings.
type WhisperEngine struct {
	ModelPath string
}

// NewWhisperEngine creates a new instance of the Whisper-based STT engine.
func NewWhisperEngine(modelPath string) *WhisperEngine {
	return &WhisperEngine{
		ModelPath: modelPath,
	}
}

// Transcribe simulates the speech-to-text process.
// It accepts an AudioChunk and returns the transcribed text.
func (e *WhisperEngine) Transcribe(chunk models.AudioChunk) (string, error) {
	// Log the transcription attempt for debugging
	log.Printf("WhisperEngine: Transcribing chunk %s", chunk.ID)

	// Simulate processing delay (e.g., 500ms)
	// This represents the time the AI would take to process the audio.
	time.Sleep(500 * time.Millisecond)

	// Return a dummy transcription for development and testing.
	return "Hola, esta es una transcripción simulada por WhisperEngine.", nil
}
