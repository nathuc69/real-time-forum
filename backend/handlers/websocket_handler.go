package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/backend/domain"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// En production, vérifier l'origine de manière plus stricte
			return true
		},
	}

	// Hub global pour gérer toutes les connexions WebSocket
	wsHub *Hub

	// Service de messages pour la persistance
	messageService domain.MessageService
)

// InitWebSocketHub initialise le hub WebSocket
func InitWebSocketHub(msgService domain.MessageService) {
	messageService = msgService
	wsHub = NewHub(messageService)
	go wsHub.Run()
	log.Println("WebSocket Hub initialized and running")
}

// WebSocketHandler gère les connexions WebSocket
func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur depuis le contexte (middleware d'authentification)
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Upgrader la connexion HTTP vers WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v", err)
		return
	}

	// Créer un nouveau client
	client := &Client{
		ID:       user.ID,
		Username: user.Username,
		Hub:      wsHub,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	// Enregistrer le client dans le hub
	client.Hub.Register <- client

	// Envoyer la liste des utilisateurs en ligne au nouveau client
	onlineUsers := wsHub.GetOnlineUsers()
	wsMsg := WSMessage{
		Type:    "online_users",
		Payload: onlineUsers,
	}

	data, err := json.Marshal(wsMsg)
	if err == nil {
		client.Send <- data
	}

	// Démarrer les goroutines pour lire et écrire
	go client.WritePump()
	go client.ReadPump()
}

// GetChatUsersHandler retourne la liste des utilisateurs avec lesquels l'utilisateur a échangé
func GetChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	chatUsers, err := messageService.GetChatList(user.ID)
	if err != nil {
		log.Printf("Error getting chat users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Ajouter le statut en ligne pour chaque utilisateur
	for i := range chatUsers {
		chatUsers[i].IsOnline = wsHub.IsUserOnline(chatUsers[i].ID)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chatUsers)
}

// GetConversationHandler retourne la conversation entre deux utilisateurs
func GetConversationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Récupérer l'ID de l'autre utilisateur depuis les paramètres de requête
	otherUserIDStr := r.URL.Query().Get("userId")
	if otherUserIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	var otherUserID int64
	if _, err := fmt.Sscanf(otherUserIDStr, "%d", &otherUserID); err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Récupérer les paramètres de pagination
	limit := 10
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		fmt.Sscanf(limitStr, "%d", &limit)
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		fmt.Sscanf(offsetStr, "%d", &offset)
	}

	messages, err := messageService.GetConversation(user.ID, otherUserID, limit, offset)
	if err != nil {
		log.Printf("Error getting conversation: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

// GetAllUsersHandler retourne tous les utilisateurs (pour pouvoir démarrer une nouvelle conversation)
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Récupérer tous les utilisateurs sauf l'utilisateur actuel
	users, err := clientService.GetAllUsers()
	if err != nil {
		log.Printf("Error getting all users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Filtrer l'utilisateur actuel et ajouter le statut en ligne
	var filteredUsers []map[string]interface{}
	for _, u := range users {
		if u.ID != user.ID {
			filteredUsers = append(filteredUsers, map[string]interface{}{
				"id":       u.ID,
				"username": u.Username,
				"isOnline": wsHub.IsUserOnline(u.ID),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filteredUsers)
}
