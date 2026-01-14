package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/domain"
	"strconv"
	"time"
)

func NewCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Extraire le postId de l'URL
	postIdStr := r.PathValue("postId")
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Récupérer l'utilisateur depuis le contexte (mis par le middleware HandleAuth)
	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := user.ID

	// Décoder le contenu du commentaire
	var requestBody struct {
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if requestBody.Content == "" {
		http.Error(w, "Comment content cannot be empty", http.StatusBadRequest)
		return
	}

	// Créer le nouveau commentaire
	newComment := &domain.Comment{
		PostID:    postId,
		UserID:    userID,
		Content:   requestBody.Content,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err = commentsService.NewCommentService(newComment)
	if err != nil {
		fmt.Println("Error creating comment:", err)
		http.Error(w, "Error creating comment", http.StatusInternalServerError)
		return
	}

	// Utiliser le username de l'utilisateur connecté
	newComment.Username = user.Username

	// Retourner le commentaire créé
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newComment)
}
