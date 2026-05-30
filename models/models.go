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

// Translation represents a translated text and its corresponding synthetic audio.
type Translation struct {
	Text  string `json:"text"`
	Audio []byte `json:"audio,omitempty"`
}

// ProcessedText represents the result of STT, translation, and TTS processes.
type ProcessedText struct {
	OriginalChunkID string                 `json:"original_chunk_id"`
	SpeakerID       string                 `json:"speaker_id"`
	OriginalLang    string                 `json:"original_lang"`
	OriginalText    string                 `json:"original_text"`
	Translations    map[string]Translation `json:"translations"`
	Timestamp       time.Time              `json:"timestamp"`
}

// RoutingRule represents a configuration rule for translations.
// It maps an incoming source language to an expected target language.
type RoutingRule struct {
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
	Active     bool   `json:"active"`
}

// ButtonEvent represents a hardware event from a PTT button or similar device.
type ButtonEvent struct {
	DeviceID      string `json:"device_id"`
	Action        string `json:"action"` // "press" or "release"
	LangRequested string `json:"lang_requested"`
}
