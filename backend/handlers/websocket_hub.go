package handlers

import (
	"encoding/json"
	"log"
	"real-time-forum/backend/domain"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client représente une connexion WebSocket d'un utilisateur
type Client struct {
	ID       int64
	Username string
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
}

// Hub maintient l'ensemble des clients actifs et diffuse les messages
type Hub struct {
	// Clients enregistrés indexés par leur ID utilisateur
	Clients map[int64]*Client

	// Messages entrants des clients
	Broadcast chan []byte

	// Enregistrer les demandes des clients
	Register chan *Client

	// Désenregistrer les demandes des clients
	Unregister chan *Client

	// Mutex pour protéger l'accès concurrent à la map des clients
	mu sync.RWMutex

	// Service de messages pour la persistance
	MessageService domain.MessageService
}

// WSMessage représente un message WebSocket
type WSMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp time.Time   `json:"timestamp"`
}

// ChatMessage représente un message de chat
type ChatMessage struct {
	SenderID       int64  `json:"senderId"`
	ReceiverID     int64  `json:"receiverId"`
	Content        string `json:"content"`
	SenderUsername string `json:"senderUsername,omitempty"`
}

// UserStatus représente le statut d'un utilisateur
type UserStatus struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	IsOnline bool   `json:"isOnline"`
}

// NewHub crée une nouvelle instance de Hub
func NewHub(messageService domain.MessageService) *Hub {
	return &Hub{
		Clients:        make(map[int64]*Client),
		Broadcast:      make(chan []byte, 256),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		MessageService: messageService,
	}
}

// Run démarre le hub et gère les événements
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()

			log.Printf("Client connected: %s (ID: %d)", client.Username, client.ID)

			// Notifier tous les autres clients que cet utilisateur est en ligne
			h.BroadcastUserStatus(client.ID, client.Username, true)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
				h.mu.Unlock()

				log.Printf("Client disconnected: %s (ID: %d)", client.Username, client.ID)

				// Notifier tous les autres clients que cet utilisateur est hors ligne
				h.BroadcastUserStatus(client.ID, client.Username, false)
			} else {
				h.mu.Unlock()
			}

		case message := <-h.Broadcast:
			// Traiter le message broadcast (si nécessaire)
			log.Printf("Broadcasting message: %s", string(message))
		}
	}
}

// SendToUser envoie un message à un utilisateur spécifique
func (h *Hub) SendToUser(userID int64, message []byte) {
	h.mu.RLock()
	client, ok := h.Clients[userID]
	h.mu.RUnlock()

	if ok {
		select {
		case client.Send <- message:
			// Message envoyé avec succès
		default:
			// Le canal est plein : ne pas fermer la connexion pour éviter de couper
			// la connexion de l'expéditeur lui-même pendant HandleChatMessage.
			// On laisse juste tomber le message et on logge un avertissement.
			log.Printf("⚠️ Send channel full for user %d, dropping message", userID)
		}
	}
}

// BroadcastUserStatus diffuse le statut d'un utilisateur à tous les clients connectés
func (h *Hub) BroadcastUserStatus(userID int64, username string, isOnline bool) {
	status := UserStatus{
		UserID:   userID,
		Username: username,
		IsOnline: isOnline,
	}

	wsMsg := WSMessage{
		Type:      "user_status",
		Payload:   status,
		Timestamp: time.Now(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshaling user status: %v", err)
		return
	}

	// Collecter les clients à notifier sous RLock (lecture seule)
	h.mu.RLock()
	targets := make([]*Client, 0)
	for id, client := range h.Clients {
		if id != userID {
			targets = append(targets, client)
		}
	}
	h.mu.RUnlock()

	// Envoyer sans tenir le lock (évite deadlock et data race)
	for _, client := range targets {
		select {
		case client.Send <- data:
		default:
			// Canal plein, on logue sans fermer la connexion
			log.Printf("⚠️ Send channel full for client %s, dropping status update", client.Username)
		}
	}
}

// GetOnlineUsers retourne la liste des utilisateurs en ligne
func (h *Hub) GetOnlineUsers() []UserStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]UserStatus, 0, len(h.Clients))
	for _, client := range h.Clients {
		users = append(users, UserStatus{
			UserID:   client.ID,
			Username: client.Username,
			IsOnline: true,
		})
	}

	return users
}

// IsUserOnline vérifie si un utilisateur est en ligne
func (h *Hub) IsUserOnline(userID int64) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.Clients[userID]
	return ok
}

// ReadPump pompe les messages de la connexion WebSocket vers le hub
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Traiter le message
		c.HandleMessage(message)
	}
}

// WritePump pompe les messages du hub vers la connexion WebSocket
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// Le hub a fermé le canal
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Ajouter les messages en file d'attente au message actuel
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HandleMessage traite les messages reçus du client
func (c *Client) HandleMessage(message []byte) {
	var wsMsg WSMessage
	if err := json.Unmarshal(message, &wsMsg); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	switch wsMsg.Type {
	case "chat_message":
		c.HandleChatMessage(wsMsg.Payload)
	case "typing":
		c.HandleTyping(wsMsg.Payload)
	case "mark_read":
		c.HandleMarkRead(wsMsg.Payload)
	default:
		log.Printf("Unknown message type: %s", wsMsg.Type)
	}
}

// HandleChatMessage traite un message de chat
func (c *Client) HandleChatMessage(payload interface{}) {
	log.Printf("📨 HandleChatMessage called for client %s (ID: %d)", c.Username, c.ID)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling chat message payload: %v", err)
		return
	}

	log.Printf("📨 Payload data: %s", string(data))

	var chatMsg ChatMessage
	if err := json.Unmarshal(data, &chatMsg); err != nil {
		log.Printf("Error unmarshaling chat message: %v", err)
		return
	}

	log.Printf("📨 Chat message: from %d to %d, content: %s", chatMsg.SenderID, chatMsg.ReceiverID, chatMsg.Content)

	// Vérifier que l'expéditeur est bien le client connecté
	if chatMsg.SenderID != c.ID {
		log.Printf("Sender ID mismatch: %d != %d", chatMsg.SenderID, c.ID)
		return
	}

	// Sauvegarder le message dans la base de données
	msg := &domain.Message{
		SenderID:   chatMsg.SenderID,
		ReceiverID: chatMsg.ReceiverID,
		Content:    chatMsg.Content,
		CreatedAt:  time.Now(),
		IsRead:     false,
	}

	log.Printf("💾 Attempting to save message to database...")
	if err := c.Hub.MessageService.SendMessage(msg); err != nil {
		log.Printf("❌ Error saving message: %v", err)
		return
	}
	log.Printf("✅ Message saved successfully with ID: %d", msg.ID)

	// Ajouter le nom d'utilisateur au message
	chatMsg.SenderUsername = c.Username

	// Créer le message WebSocket
	wsMsg := WSMessage{
		Type:      "chat_message",
		Payload:   chatMsg,
		Timestamp: msg.CreatedAt,
	}

	responseData, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		return
	}

	log.Printf("📤 Sending message to receiver %d and sender %d", chatMsg.ReceiverID, chatMsg.SenderID)

	// Envoyer le message au destinataire
	c.Hub.SendToUser(chatMsg.ReceiverID, responseData)

	// Envoyer une confirmation à l'expéditeur
	c.Hub.SendToUser(chatMsg.SenderID, responseData)
}

// HandleTyping traite un événement de frappe
func (c *Client) HandleTyping(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var typingData struct {
		ReceiverID int64 `json:"receiverId"`
	}

	if err := json.Unmarshal(data, &typingData); err != nil {
		return
	}

	wsMsg := WSMessage{
		Type: "typing",
		Payload: map[string]interface{}{
			"userId":   c.ID,
			"username": c.Username,
		},
		Timestamp: time.Now(),
	}

	responseData, _ := json.Marshal(wsMsg)
	c.Hub.SendToUser(typingData.ReceiverID, responseData)
}

// HandleMarkRead marque les messages comme lus
func (c *Client) HandleMarkRead(payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}

	var markReadData struct {
		SenderID int64 `json:"senderId"`
	}

	if err := json.Unmarshal(data, &markReadData); err != nil {
		return
	}

	// Marquer les messages comme lus
	if err := c.Hub.MessageService.MarkAsRead(markReadData.SenderID, c.ID); err != nil {
		log.Printf("Error marking messages as read: %v", err)
	}
}
