package domain

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"postId"`
	UserID    int64  `json:"userId"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	CreatedAt string `json:"createdAt"`
}

type CommentsRepository interface {
	GetCommentsByPostIDRepo(postID int64) ([]*Comment, error)
	GetUsernameByIdRepo(userID int64) (string, error)
	NewCommentRepo(comment *Comment) error
}

type CommentsService interface {
	NewCommentService(comment *Comment) error
	GetCommentsByPostIDService(postID int64) ([]*Comment, error)
}
