package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"speaking_hearts/models"
	"sync"
	"time"
)

// DualWriter writes byte data to two distinct directory paths simultaneously.
type DualWriter struct {
	PrimaryPath   string
	SecondaryPath string
}

// Write saves the data to both the primary and secondary storage locations.
func (dw *DualWriter) Write(filename string, data []byte) error {
	paths := []string{dw.PrimaryPath, dw.SecondaryPath}
	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
		fullPath := filepath.Join(path, filename)
		if err := os.WriteFile(fullPath, data, 0644); err != nil {
			return fmt.Errorf("failed to write to %s: %v", fullPath, err)
		}
	}
	return nil
}

// FragmentRecorder manages the accumulation of audio and transcriptions.
type FragmentRecorder struct {
	Writer         *DualWriter
	Interval       time.Duration
	audioBuffer    []models.AudioChunk
	textBuffer     []models.ProcessedText
	mu             sync.Mutex
	fragmentID     int
}

// NewFragmentRecorder creates a new instance of FragmentRecorder.
func NewFragmentRecorder(writer *DualWriter, interval time.Duration) *FragmentRecorder {
	return &FragmentRecorder{
		Writer:      writer,
		Interval:    interval,
		audioBuffer: make([]models.AudioChunk, 0),
		textBuffer:  make([]models.ProcessedText, 0),
	}
}

// Run starts the recording loop, capturing both audio and text data.
func (r *FragmentRecorder) Run(audioChan <-chan models.AudioChunk, textChan <-chan models.ProcessedText) {
	log.Printf("Fragment Recorder started (Interval: %v)", r.Interval)
	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()

	for {
		select {
		case chunk, ok := <-audioChan:
			if !ok {
				r.flush()
				return
			}
			r.mu.Lock()
			r.audioBuffer = append(r.audioBuffer, chunk)
			r.mu.Unlock()

		case text, ok := <-textChan:
			if !ok {
				continue
			}
			r.mu.Lock()
			r.textBuffer = append(r.textBuffer, text)
			r.mu.Unlock()

		case <-ticker.C:
			r.flush()
		}
	}
}

// flush writes the buffered data to disk and resets the state.
func (r *FragmentRecorder) flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.audioBuffer) == 0 && len(r.textBuffer) == 0 {
		return
	}

	r.fragmentID++
	timestamp := time.Now().Format("20060102_150405")
	baseName := fmt.Sprintf("fragment_%d_%s", r.fragmentID, timestamp)
	
	// Save Mock Audio
	mockAudioData := []byte(fmt.Sprintf("MOCK WAV DATA for fragment %d", r.fragmentID))
	audioFile := baseName + ".wav"
	if err := r.Writer.Write(audioFile, mockAudioData); err != nil {
		log.Printf("Recorder Error: %v", err)
	}

	// Generate and save complete metadata
	metadata := struct {
		FragmentID     int                    `json:"fragment_id"`
		Timestamp      time.Time              `json:"timestamp"`
		ChunkCount     int                    `json:"chunk_count"`
		Chunks         []models.AudioChunk    `json:"chunks"`
		Transcriptions []models.ProcessedText `json:"transcriptions"`
	}{
		FragmentID:     r.fragmentID,
		Timestamp:      time.Now(),
		ChunkCount:     len(r.audioBuffer),
		Chunks:         r.audioBuffer,
		Transcriptions: r.textBuffer,
	}

	metaJSON, _ := json.MarshalIndent(metadata, "", "  ")
	metaFile := baseName + ".json"
	if err := r.Writer.Write(metaFile, metaJSON); err != nil {
		log.Printf("Recorder Error: %v", err)
	}

	log.Printf("Recorder: Flushed fragment %d (Chunks: %d, Text: %d)", 
		r.fragmentID, len(r.audioBuffer), len(r.textBuffer))
	
	r.audioBuffer = make([]models.AudioChunk, 0)
	r.textBuffer = make([]models.ProcessedText, 0)
}
