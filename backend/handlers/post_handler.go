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

	w.Header().Set("Content-Type", "application/json")

	posts, err := postsService.GetAllPostsService()
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

	post, err := postsService.GetPostByIDService(postID)
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
