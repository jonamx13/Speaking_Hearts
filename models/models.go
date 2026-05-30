package models

import (
	"time"
)

// AudioChunk represents a raw segment of audio data captured from the hardware.
type AudioChunk struct {
	ID         string    `json:"id"`
	Source     string    `json:"source"`
	Timestamp  time.Time `json:"timestamp"`
	Data       []float32 `json:"data"`
	SampleRate int       `json:"sample_rate"`
	LangIn     string    `json:"lang_in"`
}

// ProcessedText represents the result of STT and translation processes.
type ProcessedText struct {
	OriginalChunkID string            `json:"original_chunk_id"`
	SpeakerID       string            `json:"speaker_id"`
	OriginalLang    string            `json:"original_lang"`
	OriginalText    string            `json:"original_text"`
	Translations    map[string]string `json:"translations"`
	Timestamp       time.Time         `json:"timestamp"`
}

// RoutingRule represents a configuration rule for translations.
// It maps an incoming source language to an expected target language.
type RoutingRule struct {
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	Active     bool   `json:"active"`
}
