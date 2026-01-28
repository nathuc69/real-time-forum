package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/domain"
	"strconv"
	"strings"
)

// GetUserReactionHandler récupère la réaction de l'utilisateur pour un post
func GetUserReactionHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur depuis le contexte
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID du post depuis l'URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postIDStr := pathParts[3]
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	result, err := postsService.GetUserReactionService(postID, user.ID)
	if err != nil {
		http.Error(w, "Error retrieving reaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// AddReactionHandler ajoute ou modifie une réaction
func AddReactionHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur depuis le contexte
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID du post depuis l'URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postIDStr := pathParts[3]
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	// Décoder le body de la requête
	var requestBody struct {
		Value int `json:"value"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Valider la valeur (1 = like, -1 = dislike)
	if requestBody.Value != 1 && requestBody.Value != -1 {
		http.Error(w, "Invalid reaction value. Must be 1 or -1", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	result, err := postsService.AddReactionService(postID, user.ID, requestBody.Value)
	if err != nil {
		http.Error(w, "Error adding reaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// RemoveReactionHandler supprime une réaction
func RemoveReactionHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'utilisateur depuis le contexte
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extraire l'ID du post depuis l'URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	postIDStr := pathParts[3]
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID format", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	result, err := postsService.RemoveReactionService(postID, user.ID)
	if err != nil {
		http.Error(w, "Error removing reaction", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

// ReactionHandler gère les différentes méthodes HTTP pour les réactions
func ReactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetUserReactionHandler(w, r)
	case http.MethodPost:
		AddReactionHandler(w, r)
	case http.MethodDelete:
		RemoveReactionHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
