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
// This provides a redundancy layer for the ceremonial recordings.
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

// FragmentRecorder manages the accumulation and periodic flushing of audio data.
type FragmentRecorder struct {
	Writer     *DualWriter
	Interval   time.Duration
	buffer     []models.AudioChunk
	mu         sync.Mutex
	fragmentID int
}

// NewFragmentRecorder creates a new instance of FragmentRecorder.
func NewFragmentRecorder(writer *DualWriter, interval time.Duration) *FragmentRecorder {
	return &FragmentRecorder{
		Writer:   writer,
		Interval: interval,
		buffer:   make([]models.AudioChunk, 0),
	}
}

// Run starts the recording loop, flushing data at the specified interval.
func (r *FragmentRecorder) Run(audioChan <-chan models.AudioChunk) {
	log.Printf("Fragment Recorder started (Interval: %v)", r.Interval)
	ticker := time.NewTicker(r.Interval)
	defer ticker.Stop()

	for {
		select {
		case chunk, ok := <-audioChan:
			if !ok {
				r.flush() // Final flush before exiting
				return
			}
			r.mu.Lock()
			r.buffer = append(r.buffer, chunk)
			r.mu.Unlock()

		case <-ticker.C:
			r.flush()
		}
	}
}

// flush writes the buffered chunks to disk and clears the buffer.
func (r *FragmentRecorder) flush() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.buffer) == 0 {
		return
	}

	r.fragmentID++
	timestamp := time.Now().Format("20060102_150405")
	baseName := fmt.Sprintf("fragment_%d_%s", r.fragmentID, timestamp)
	
	// Prepare audio data (Mocking .wav content for now)
	// In production, this would involve proper WAV header construction and PCM encoding.
	mockAudioData := []byte(fmt.Sprintf("MOCK WAV DATA for fragment %d containing %d chunks", r.fragmentID, len(r.buffer)))
	
	audioFile := baseName + ".wav"
	if err := r.Writer.Write(audioFile, mockAudioData); err != nil {
		log.Printf("Recorder Error: Failed to write audio fragment: %v", err)
	}

	// Generate and save metadata
	metadata := struct {
		FragmentID int                 `json:"fragment_id"`
		Timestamp  time.Time           `json:"timestamp"`
		ChunkCount int                 `json:"chunk_count"`
		Chunks     []models.AudioChunk `json:"chunks"`
	}{
		FragmentID: r.fragmentID,
		Timestamp:  time.Now(),
		ChunkCount: len(r.buffer),
		Chunks:     r.buffer,
	}

	metaJSON, _ := json.MarshalIndent(metadata, "", "  ")
	metaFile := baseName + ".json"
	if err := r.Writer.Write(metaFile, metaJSON); err != nil {
		log.Printf("Recorder Error: Failed to write metadata: %v", err)
	}

	log.Printf("Recorder: Flushed fragment %d with %d chunks", r.fragmentID, len(r.buffer))
	
	// Reset buffer for the next fragment
	r.buffer = make([]models.AudioChunk, 0)
}
