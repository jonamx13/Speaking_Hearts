package broadcast

import (
	"encoding/json"
	"log"
	"speaking_hearts/models"
)

// Manager maintains the set of active clients and broadcasts messages to them.
// It acts as the central hub for all WebSocket communication.
type Manager struct {
	// Registered clients.
	Clients map[*Client]bool

	// Inbound messages from the processing layers.
	Broadcast chan models.ProcessedText

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client
}

// NewManager creates a new instance of the Broadcast Manager.
func NewManager() *Manager {
	return &Manager{
		Broadcast:  make(chan models.ProcessedText),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
	}
}

// Run is the main loop for the Broadcast Manager, handling concurrency through channels.
// This design avoids the need for explicit mutexes by serializing all state changes.
func (m *Manager) Run() {
	log.Println("Broadcast Manager started")
	for {
		select {
		case client := <-m.Register:
			// Add a new client to the registry.
			m.Clients[client] = true
			log.Printf("Client %s registered (Lang: %s)", client.ID, client.PreferredLang)

		case client := <-m.Unregister:
			// Remove a client and close its send channel.
			if _, ok := m.Clients[client]; ok {
				delete(m.Clients, client)
				close(client.Send)
				log.Printf("Client %s unregistered", client.ID)
			}

		case processedText := <-m.Broadcast:
			// Broadcast the processed text to all clients.
			// Each client receives only the translation for their PreferredLang.
			for client := range m.Clients {
				// Prepare the message payload specifically for this client.
				payload := struct {
					OriginalChunkID string            `json:"original_chunk_id"`
					SpeakerID       string            `json:"speaker_id"`
					Text            string            `json:"text"`
					Timestamp       string            `json:"timestamp"`
				}{
					OriginalChunkID: processedText.OriginalChunkID,
					SpeakerID:       processedText.SpeakerID,
					Timestamp:       processedText.Timestamp.Format("15:04:05"),
				}

				// Check if a translation for the client's preferred language exists.
				// If not, default to the original text.
				if translation, exists := processedText.Translations[client.PreferredLang]; exists {
					payload.Text = translation
				} else {
					payload.Text = processedText.OriginalText
				}

				// Serialize to JSON and send.
				msg, err := json.Marshal(payload)
				if err != nil {
					log.Printf("Error marshaling broadcast message: %v", err)
					continue
				}

				select {
				case client.Send <- msg:
					// Message queued successfully.
				default:
					// If the client's send buffer is full, drop the client.
					close(client.Send)
					delete(m.Clients, client)
				}
			}
		}
	}
}
