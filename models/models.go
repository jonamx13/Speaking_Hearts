package models

import (
	"time"
)

// AudioChunk represents a raw segment of audio data captured from the hardware.
type AudioChunk struct {
	ID         string
	Source     string
	Timestamp  time.Time
	Data       []float32
	SampleRate int
	LangIn     string
}

// ProcessedText represents the result of STT and translation processes.
type ProcessedText struct {
	OriginalChunkID string
	SpeakerID       string
	OriginalLang    string
	OriginalText    string
	Translations    map[string]string
	Timestamp       time.Time
}
