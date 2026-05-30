package tts

import (
	"log"
	"time"
)

// Engine defines the interface for text-to-speech generation.
type Engine interface {
	GenerateSpeech(text string, lang string) ([]byte, error)
}

// MockEngine simulates a TTS service for development and testing.
type MockEngine struct{}

// NewMockEngine creates a new instance of the MockEngine.
func NewMockEngine() *MockEngine {
	return &MockEngine{}
}

// GenerateSpeech simulates speech generation with a delay.
func (e *MockEngine) GenerateSpeech(text string, lang string) ([]byte, error) {
	log.Printf("TTS Engine: Generating speech for [%s]: %s", lang, text)

	// Simulate processing delay (e.g., 200ms)
	time.Sleep(200 * time.Millisecond)

	// Return a mock byte array representing audio data
	return []byte("mock_audio_data_" + lang), nil
}
