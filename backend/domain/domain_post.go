package domain

type Post struct {
	ID              int64      `json:"id"`
	Title           string     `json:"title"`
	Content         string     `json:"content"`
	UserID          int64      `json:"userId"`
	Username        string     `json:"username"`
	Likes           int        `json:"likes"`
	Dislikes        int        `json:"dislikes"`
	IslikeOrDislike int        `json:"islikeordislike"`
	Comments        int        `json:"comments"`
	CreatedAt       string     `json:"createdAt"`
	CommentsList    []*Comment `json:"commentsList,omitempty"`
}

type PostsRepository interface {
	GetAllPostsRepo() ([]*Post, error)
	GetPostByIDRepo(postID int64) (*Post, error)
	GetUsernameByIdRepo(userID int64) (string, error)
	CountLikesByPostID(postID int64) (int, error)
	CountCommentsByPostID(postID int64) (int, error)
	CountDislikesByPostID(postID int64) (int, error)
	LikeOrDislikePostRepo(postID int64, userID int64) (int, error)
}

type PostsService interface {
	GetAllPostsService() ([]*Post, error)
	GetPostByIDService(postID int64, userId int64) (*Post, error)
}
