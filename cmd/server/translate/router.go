package translate

import (
	"log"
	"speaking_hearts/models"
	"sync"
)

// LanguageRouter handles the logic of determining which translations are needed
// and orchestrating the translation process.
type LanguageRouter struct {
	Translator Translator
	// Rules defines the list of language routing rules.
	// In a more dynamic system, this would be updated based on active clients.
	Rules []models.RoutingRule
	mu    sync.RWMutex
}

// NewLanguageRouter creates a new instance of the LanguageRouter.
func NewLanguageRouter(translator Translator, rules []models.RoutingRule) *LanguageRouter {
	return &LanguageRouter{
		Translator: translator,
		Rules:      rules,
	}
}

// SetRules updates the routing rules in a thread-safe manner.
func (r *LanguageRouter) SetRules(rules []models.RoutingRule) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Rules = rules
}

// RouteProcess populates the Translations map in ProcessedText by calling the Translator.
// This function implements the routing logic specified in Phase 3.
func (r *LanguageRouter) RouteProcess(p *models.ProcessedText) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if p.Translations == nil {
		p.Translations = make(map[string]models.Translation)
	}

	for _, rule := range r.Rules {
		// Skip if the rule is not active
		if !rule.Active {
			continue
		}

		// Apply rule only if it matches the original language
		if rule.SourceLang != "" && rule.SourceLang != p.OriginalLang {
			continue
		}

		targetLang := rule.TargetLang

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
		p.Translations[targetLang] = models.Translation{
			Text: translatedText,
		}
	}
}
