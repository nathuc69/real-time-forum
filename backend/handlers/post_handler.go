package handlers

import (
	"encoding/json"
	"net/http"
	"real-time-forum/backend/domain"
	"strconv"
	"strings"
)

var postsService domain.PostsService
var commentsService domain.CommentsService

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Récupérer l'utilisateur depuis le contexte (optionnel)
	var userID int64 = 0
	user, ok := r.Context().Value("user").(*domain.User)
	if ok && user != nil {
		userID = user.ID
	}

	w.Header().Set("Content-Type", "application/json")

	sortBy := r.URL.Query().Get("sort")
	category := r.URL.Query().Get("category")

	posts, err := postsService.GetAllPostsService(userID, sortBy, category)
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetPostByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Authentification optionnelle
	var userID int64 = 0
	user, ok := r.Context().Value("user").(*domain.User)
	if ok && user != nil {
		userID = user.ID
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

	post, err := postsService.GetPostByIDService(postID, userID)
	if err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Récupérer les commentaires seulement si le service est disponible
	if commentsService != nil {
		comments, err := commentsService.GetCommentsByPostIDService(postID)
		if err == nil {
			post.CommentsList = comments
		}
	}

	err = json.NewEncoder(w).Encode(post)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, ok := r.Context().Value("user").(*domain.User)
	if !ok || user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Title      string   `json:"title"`
		Content    string   `json:"content"`
		Categories []string `json:"categories"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	post := &domain.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  user.ID,
	}

	postID, err := postsService.CreatePostService(post, req.Categories)
	if err != nil {
		http.Error(w, "Error creating post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      postID,
		"message": "Post created successfully",
	})
}

func GetAllCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	categories, err := postsService.GetAllCategoriesService()
	if err != nil {
		http.Error(w, "Error retrieving categories", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(categories)
}
