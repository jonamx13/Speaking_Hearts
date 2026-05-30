package stt

import (
	"log"
	"speaking_hearts/models"
	"sync"
	"time"
)

// WorkerPool manages a pool of concurrent STT workers.
// This pattern ensures that we can scale our processing capacity to handle
// high-frequency audio chunks without blocking the acquisition layer.
type WorkerPool struct {
	WorkerCount int
	InputChan   <-chan models.AudioChunk
	OutputChan  chan<- models.ProcessedText
	Engine      STTEngine
	wg          sync.WaitGroup
	quit        chan struct{}
}

// NewWorkerPool creates a new instance of the STT WorkerPool.
func NewWorkerPool(count int, input <-chan models.AudioChunk, output chan<- models.ProcessedText, engine STTEngine) *WorkerPool {
	return &WorkerPool{
		WorkerCount: count,
		InputChan:   input,
		OutputChan:  output,
		Engine:      engine,
		quit:        make(chan struct{}),
	}
}

// Start initializes the worker goroutines.
func (p *WorkerPool) Start() {
	log.Printf("Starting STT Worker Pool with %d workers", p.WorkerCount)
	p.wg.Add(p.WorkerCount)

	for i := 0; i < p.WorkerCount; i++ {
		go p.worker(i)
	}
}

// Stop signals all workers to stop and waits for them to finish.
func (p *WorkerPool) Stop() {
	log.Println("Stopping STT Worker Pool...")
	close(p.quit)
	p.wg.Wait()
	log.Println("STT Worker Pool stopped gracefully")
}

// worker is the internal loop for each processing goroutine.
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	log.Printf("STT Worker %d started", id)

	for {
		select {
		case chunk, ok := <-p.InputChan:
			if !ok {
				// Input channel closed, worker should exit.
				log.Printf("STT Worker %d: Input channel closed", id)
				return
			}

			// Call the polymorphic STT engine
			text, err := p.Engine.Transcribe(chunk)
			if err != nil {
				log.Printf("STT Worker %d: Error transcribing chunk %s: %v", id, chunk.ID, err)
				continue
			}

			processed := models.ProcessedText{
				OriginalChunkID: chunk.ID,
				SpeakerID:       "Worker_" + string(rune(48+id)),
				OriginalLang:    chunk.LangIn,
				OriginalText:    text,
				Translations:    make(map[string]models.Translation),
				Timestamp:       time.Now(),
			}

			// Send the result to the next layer (Translator/Broadcast)
			p.OutputChan <- processed

		case <-p.quit:
			// Quit signal received, worker should exit.
			log.Printf("STT Worker %d: Received quit signal", id)
			return
		}
	}
}
