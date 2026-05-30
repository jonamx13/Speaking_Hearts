package tts

import (
	"log"
	"speaking_hearts/models"
	"sync"
)

// WorkerPool manages concurrent TTS workers.
type WorkerPool struct {
	WorkerCount int
	InputChan   <-chan models.ProcessedText
	OutputChan  chan<- models.ProcessedText
	Engine      Engine
	wg          sync.WaitGroup
	quit        chan struct{}
}

// NewWorkerPool creates a new instance of the TTS WorkerPool.
func NewWorkerPool(count int, input <-chan models.ProcessedText, output chan<- models.ProcessedText, engine Engine) *WorkerPool {
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
	log.Printf("Starting TTS Worker Pool with %d workers", p.WorkerCount)
	p.wg.Add(p.WorkerCount)

	for i := 0; i < p.WorkerCount; i++ {
		go p.worker(i)
	}
}

// Stop signals all workers to stop.
func (p *WorkerPool) Stop() {
	close(p.quit)
	p.wg.Wait()
}

// worker processes TTS tasks concurrently.
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()
	log.Printf("TTS Worker %d started", id)

	for {
		select {
		case processed, ok := <-p.InputChan:
			if !ok {
				return
			}

			// Generate audio for each translation
			for lang, trans := range processed.Translations {
				audio, err := p.Engine.GenerateSpeech(trans.Text, lang)
				if err != nil {
					log.Printf("TTS Worker %d Error: %v", id, err)
					continue
				}

				// Update the translation with the generated audio
				trans.Audio = audio
				processed.Translations[lang] = trans
			}

			// Forward the fully complete payload (text + audio) to the next stage
			p.OutputChan <- processed

		case <-p.quit:
			return
		}
	}
}
