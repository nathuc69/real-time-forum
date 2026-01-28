package repositories

import (
	"database/sql"
	"real-time-forum/backend/domain"
)

type PostsRepo struct {
	db *sql.DB
}

func NewPostsRepository(db *sql.DB) domain.PostsRepository {
	return &PostsRepo{db: db}
}

func (r *PostsRepo) GetAllPostsRepo() ([]*domain.Post, error) {
	rows, err := r.db.Query(`
		SELECT id, title, content, user_id, created_at
		FROM posts
		ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostsRepo) GetPostByIDRepo(postID int64) (*domain.Post, error) {
	post := &domain.Post{}
	err := r.db.QueryRow(`
		SELECT id, title, content, user_id, created_at
		FROM posts
		WHERE id = ?`, postID).Scan(&post.ID, &post.Title, &post.Content, &post.UserID, &post.CreatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostsRepo) GetUsernameByIdRepo(userID int64) (string, error) {
	var username string
	err := r.db.QueryRow(`SELECT userName FROM users WHERE id = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func (r *PostsRepo) CountLikesByPostID(postID int64) (int, error) {
	var likesCount int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE target_type = 'posts' AND target_id = ? AND value = 1`, postID).Scan(&likesCount)
	if err != nil {
		return 0, err
	}
	return likesCount, nil
}
func (r *PostsRepo) CountDislikesByPostID(postID int64) (int, error) {
	var dislikesCount int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM reactions WHERE target_type = 'posts' AND target_id = ? AND value = -1`, postID).Scan(&dislikesCount)
	if err != nil {
		return 0, err
	}
	return dislikesCount, nil
}

func (r *PostsRepo) CountCommentsByPostID(postID int64) (int, error) {
	var commentsCount int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM comments WHERE post_id = ?`, postID).Scan(&commentsCount)
	if err != nil {
		return 0, err
	}
	return commentsCount, nil
}

func (r *PostsRepo) LikeOrDislikePostRepo(postID int64, userID int64) (int, error) {
	var LikeorDislike int
	err := r.db.QueryRow(`SELECT value FROM reactions WHERE target_type = 'posts' AND target_id = ? AND user_id = ?`, postID, userID).Scan(&LikeorDislike)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return LikeorDislike, nil
}

// AddReactionRepo ajoute ou met à jour une réaction
func (r *PostsRepo) AddReactionRepo(postID int64, userID int64, value int) error {
	// Vérifier si une réaction existe déjà
	var existingValue int
	err := r.db.QueryRow(`SELECT value FROM reactions WHERE target_type = 'posts' AND target_id = ? AND user_id = ?`, postID, userID).Scan(&existingValue)

	if err == sql.ErrNoRows {
		// Aucune réaction existante, en créer une nouvelle
		_, err = r.db.Exec(`INSERT INTO reactions (value, target_type, target_id, user_id) VALUES (?, 'posts', ?, ?)`, value, postID, userID)
		return err
	} else if err != nil {
		return err
	}

	// Une réaction existe déjà, la mettre à jour
	_, err = r.db.Exec(`UPDATE reactions SET value = ?, updated_at = CURRENT_TIMESTAMP WHERE target_type = 'posts' AND target_id = ? AND user_id = ?`, value, postID, userID)
	return err
}

// RemoveReactionRepo supprime une réaction
func (r *PostsRepo) RemoveReactionRepo(postID int64, userID int64) error {
	_, err := r.db.Exec(`DELETE FROM reactions WHERE target_type = 'posts' AND target_id = ? AND user_id = ?`, postID, userID)
	return err
}

// GetUserReactionRepo récupère la réaction d'un utilisateur pour un post
func (r *PostsRepo) GetUserReactionRepo(postID int64, userID int64) (int, error) {
	var value int
	err := r.db.QueryRow(`SELECT value FROM reactions WHERE target_type = 'posts' AND target_id = ? AND user_id = ?`, postID, userID).Scan(&value)
	if err == sql.ErrNoRows {
		return 0, nil // Pas de réaction trouvée
	}
	return value, err
}
