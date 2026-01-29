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
	UserReaction    *int       `json:"userReaction,omitempty"`
	Comments        int        `json:"comments"`
	CreatedAt       string     `json:"createdAt"`
	CommentsList    []*Comment `json:"commentsList,omitempty"`
	Categories      []string   `json:"categories,omitempty"`
}

type PostsRepository interface {
	GetAllPostsRepo(sortBy string, category string) ([]*Post, error)
	GetPostByIDRepo(postID int64) (*Post, error)
	GetUsernameByIdRepo(userID int64) (string, error)
	CountLikesByPostID(postID int64) (int, error)
	CountCommentsByPostID(postID int64) (int, error)
	CountDislikesByPostID(postID int64) (int, error)
	LikeOrDislikePostRepo(postID int64, userID int64) (int, error)
	AddReactionRepo(postID int64, userID int64, value int) error
	RemoveReactionRepo(postID int64, userID int64) error
	GetUserReactionRepo(postID int64, userID int64) (int, error)
	CreatePostRepo(post *Post) (int64, error)
	AddCategoryToPostRepo(postID int64, categoryID int64) error
	GetCategoryByNameRepo(name string) (int64, error)
	CreateCategoryRepo(name string) (int64, error)
	GetCategoriesByPostIDRepo(postID int64) ([]string, error)
	GetAllCategoriesRepo() ([]string, error)
}

type PostsService interface {
	GetAllPostsService(userId int64, sortBy string, category string) ([]*Post, error)
	GetPostByIDService(postID int64, userId int64) (*Post, error)
	AddReactionService(postID int64, userID int64, value int) (map[string]interface{}, error)
	RemoveReactionService(postID int64, userID int64) (map[string]interface{}, error)
	GetUserReactionService(postID int64, userID int64) (map[string]interface{}, error)
	CreatePostService(post *Post, categories []string) (int64, error)
	GetAllCategoriesService() ([]string, error)
}
