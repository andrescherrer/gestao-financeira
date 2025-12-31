package websocket

import (
	"sync"

	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients mapped by user ID
	clients map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client registration/unregistration and message broadcasting.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.userID] == nil {
				h.clients[client.userID] = make(map[*Client]bool)
			}
			h.clients[client.userID][client] = true
			h.mu.Unlock()

			log.Info().
				Str("user_id", client.userID).
				Int("total_clients", len(h.clients[client.userID])).
				Msg("WebSocket client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.userID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.send)
					if len(clients) == 0 {
						delete(h.clients, client.userID)
					}
				}
			}
			h.mu.Unlock()

			log.Info().
				Str("user_id", client.userID).
				Msg("WebSocket client unregistered")

		case <-h.broadcast:
			// Broadcast message to all clients (if needed)
			// For now, we'll use SendToUser for targeted messages
			log.Debug().Msg("Broadcast message received (not implemented for targeted notifications)")
		}
	}
}

// SendToUser sends a message to all clients of a specific user.
func (h *Hub) SendToUser(userID string, message []byte) {
	h.mu.RLock()
	clients, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		log.Debug().Str("user_id", userID).Msg("No WebSocket clients found for user")
		return
	}

	h.mu.RLock()
	for client := range clients {
		select {
		case client.send <- message:
		default:
			// Client's send buffer is full, close the connection
			close(client.send)
			delete(clients, client)
		}
	}
	h.mu.RUnlock()

	log.Debug().
		Str("user_id", userID).
		Int("clients_notified", len(clients)).
		Msg("Message sent to user via WebSocket")
}

// GetClientCount returns the number of active clients for a user.
func (h *Hub) GetClientCount(userID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID])
}

// GetTotalClients returns the total number of active clients across all users.
func (h *Hub) GetTotalClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, clients := range h.clients {
		total += len(clients)
	}
	return total
}
