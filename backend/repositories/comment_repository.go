package repositories

import (
	"database/sql"
	"real-time-forum/backend/domain"
)

type CommentsRepo struct {
	db *sql.DB
}

func NewCommentsRepository(db *sql.DB) domain.CommentsRepository {
	return &CommentsRepo{db: db}
}

func (r *CommentsRepo) GetCommentsByPostIDRepo(postID int64) ([]*domain.Comment, error) {
	rows, err := r.db.Query(`
		SELECT id, post_id, user_id, content, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at ASC`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		comment := &domain.Comment{}
		err := rows.Scan(&comment.ID, &comment.PostID, &comment.UserID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}
func (r *CommentsRepo) GetUsernameByIdRepo(userID int64) (string, error) {
	var username string
	err := r.db.QueryRow(`SELECT userName FROM users WHERE id = ?`, userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}
func (r *CommentsRepo) NewCommentRepo(comment *domain.Comment) error {
	_, err := r.db.Exec(`
		INSERT INTO comments (post_id, user_id, content, created_at)
		VALUES (?, ?, ?, datetime('now'))`,
		comment.PostID, comment.UserID, comment.Content)
	return err
}
