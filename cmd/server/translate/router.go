package translate

import (
	"log"
	"speaking_hearts/models"
)

// LanguageRouter handles the logic of determining which translations are needed
// and orchestrating the translation process.
type LanguageRouter struct {
	Translator Translator
	// Rules defines the list of languages that should always be translated.
	// In a more dynamic system, this would be updated based on active clients.
	Rules []string
}

// NewLanguageRouter creates a new instance of the LanguageRouter.
func NewLanguageRouter(translator Translator, languages []string) *LanguageRouter {
	return &LanguageRouter{
		Translator: translator,
		Rules:      languages,
	}
}

// RouteProcess populates the Translations map in ProcessedText by calling the Translator.
// This function implements the routing logic specified in Phase 3.
func (r *LanguageRouter) RouteProcess(p *models.ProcessedText) {
	if p.Translations == nil {
		p.Translations = make(map[string]string)
	}

	for _, targetLang := range r.Rules {
		// Skip if it's the same as the original language
		if targetLang == p.OriginalLang {
			continue
		}

		// Perform translation
		translatedText, err := r.Translator.Translate(p.OriginalText, p.OriginalLang, targetLang)
		if err != nil {
			log.Printf("Router Error: Failed to translate to %s: %v", targetLang, err)
			continue
		}

		// Update the struct's translation map
		p.Translations[targetLang] = translatedText
	}
}
