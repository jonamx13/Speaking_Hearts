package translate

import (
	"fmt"
	"log"
	"time"
)

// Translator defines the interface for text translation services.
type Translator interface {
	Translate(text string, sourceLang string, targetLang string) (string, error)
}

// TranslatorService is a wrapper for the NLLB-200 translation engine.
// In this phase, it acts as a mock to simulate the translation process.
type TranslatorService struct {
	ModelPath string
}

// NewTranslatorService creates a new instance of the TranslatorService.
func NewTranslatorService(modelPath string) *TranslatorService {
	return &TranslatorService{
		ModelPath: modelPath,
	}
}

// Translate simulates a translation using NLLB-200.
// It accepts source text and language codes, returning the translated string.
func (s *TranslatorService) Translate(text string, sourceLang string, targetLang string) (string, error) {
	// If source and target are the same, return original
	if sourceLang == targetLang {
		return text, nil
	}

	log.Printf("Translator: Translating from %s to %s", sourceLang, targetLang)

	// Simulate high-latency AI processing (e.g., 300ms)
	time.Sleep(300 * time.Millisecond)

	// Return a dummy translation for development and testing.
	return fmt.Sprintf("[Translated to %s]: %s", targetLang, text), nil
}
