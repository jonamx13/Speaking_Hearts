package buttons

import (
	"log"
	"speaking_hearts/cmd/server/translate"
	"speaking_hearts/models"
)

// ButtonManager handles PTT hardware events and updates routing rules.
type ButtonManager struct {
	EventChan chan models.ButtonEvent
	Router    *translate.LanguageRouter
}

// NewButtonManager creates a new instance of the ButtonManager.
func NewButtonManager(router *translate.LanguageRouter) *ButtonManager {
	return &ButtonManager{
		EventChan: make(chan models.ButtonEvent, 10),
		Router:    router,
	}
}

// Run starts the button event listener loop.
func (m *ButtonManager) Run() {
	log.Println("Button Manager started")
	for event := range m.EventChan {
		log.Printf("Button Event Received: Device=%s, Action=%s, Lang=%s", 
			event.DeviceID, event.Action, event.LangRequested)

		if event.Action == "press" {
			// Dynamic Routing Logic:
			// When a button is pressed, we update the router with new rules.
			// Example: Priest (Spanish) -> English, Russian, German
			var newRules []models.RoutingRule
			
			switch event.LangRequested {
			case "es":
				newRules = []models.RoutingRule{
					{SourceLang: "es", TargetLang: "en", Active: true},
					{SourceLang: "es", TargetLang: "ru", Active: true},
					{SourceLang: "es", TargetLang: "de", Active: true},
				}
			case "en":
				newRules = []models.RoutingRule{
					{SourceLang: "en", TargetLang: "es", Active: true},
					{SourceLang: "en", TargetLang: "ru", Active: true},
				}
			default:
				log.Printf("Warning: No specific routing rules for language %s", event.LangRequested)
				continue
			}

			log.Printf("Updating routing rules for %s speaker", event.LangRequested)
			m.Router.SetRules(newRules)
		}
		// On "release", we could either keep the rules or reset them.
		// For now, we'll keep them active until the next "press".
	}
}
