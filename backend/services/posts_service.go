package services

import (
	"errors"
	"real-time-forum/backend/domain"
	"sort"
)

var ErrInvalidInput = errors.New("invalid input")

type postsService struct {
	postsRepo domain.PostsRepository
}

func NewPostsService(pr domain.PostsRepository) domain.PostsService {
	return &postsService{postsRepo: pr}
}
func (s *postsService) GetAllPostsService(userId int64, sortBy string, category string) ([]*domain.Post, error) {
	posts, err := s.postsRepo.GetAllPostsRepo(sortBy, category)
	if err != nil {
		return nil, err
	}

	// Récupérer les noms d'utilisateur pour chaque post
	for _, post := range posts {
		username, err := s.postsRepo.GetUsernameByIdRepo(post.UserID)
		if err != nil {
			return nil, err
		}
		likesCount, err := s.postsRepo.CountLikesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		dislikesCount, err := s.postsRepo.CountDislikesByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		commentsCount, err := s.postsRepo.CountCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}

		// Récupérer la réaction de l'utilisateur si userId > 0
		var userReaction *int
		if userId > 0 {
			reaction, err := s.postsRepo.GetUserReactionRepo(post.ID, userId)
			if err != nil {
				return nil, err
			}
			if reaction != 0 {
				userReaction = &reaction
			}
		}

		categories, err := s.postsRepo.GetCategoriesByPostIDRepo(post.ID)
		if err != nil {
			return nil, err
		}

		post.Likes = likesCount
		post.Dislikes = dislikesCount
		post.Comments = commentsCount
		post.Username = username
		post.UserReaction = userReaction
		post.Categories = categories
	}

	if sortBy == "likes" {
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Likes > posts[j].Likes
		})
	}

	return posts, nil
}

func (s *postsService) GetPostByIDService(postID int64, userID int64) (*domain.Post, error) {
	post, err := s.postsRepo.GetPostByIDRepo(postID)
	if err != nil {
		return nil, err
	}

	// Récupérer les informations complémentaires
	username, err := s.postsRepo.GetUsernameByIdRepo(post.UserID)
	if err != nil {
		return nil, err
	}
	likesCount, err := s.postsRepo.CountLikesByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	dislikesCount, err := s.postsRepo.CountDislikesByPostID(post.ID)
	if err != nil {
		return nil, err
	}
	commentsCount, err := s.postsRepo.CountCommentsByPostID(post.ID)
	if err != nil {
		return nil, err
	}

	// Récupérer la réaction de l'utilisateur si userID > 0
	var userReaction *int
	if userID > 0 {
		reaction, err := s.postsRepo.GetUserReactionRepo(post.ID, userID)
		if err != nil {
			return nil, err
		}
		if reaction != 0 {
			userReaction = &reaction
		}
	}

	categories, err := s.postsRepo.GetCategoriesByPostIDRepo(post.ID)
	if err != nil {
		return nil, err
	}

	post.Dislikes = dislikesCount
	post.Likes = likesCount
	post.Comments = commentsCount
	post.Username = username
	post.UserReaction = userReaction
	post.Categories = categories
	return post, nil
}

// AddReactionService ajoute ou modifie une réaction
func (s *postsService) AddReactionService(postID int64, userID int64, value int) (map[string]interface{}, error) {
	// Valider la valeur (doit être 1 ou -1)
	if value != 1 && value != -1 {
		return nil, ErrInvalidInput
	}

	// Ajouter ou mettre à jour la réaction
	err := s.postsRepo.AddReactionRepo(postID, userID, value)
	if err != nil {
		return nil, err
	}

	// Récupérer les nouveaux totaux
	likes, err := s.postsRepo.CountLikesByPostID(postID)
	if err != nil {
		return nil, err
	}
	dislikes, err := s.postsRepo.CountDislikesByPostID(postID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"totalLikes":    likes,
		"totalDislikes": dislikes,
		"userReaction":  value,
	}, nil
}

// RemoveReactionService supprime une réaction
func (s *postsService) RemoveReactionService(postID int64, userID int64) (map[string]interface{}, error) {
	err := s.postsRepo.RemoveReactionRepo(postID, userID)
	if err != nil {
		return nil, err
	}

	// Récupérer les nouveaux totaux
	likes, err := s.postsRepo.CountLikesByPostID(postID)
	if err != nil {
		return nil, err
	}
	dislikes, err := s.postsRepo.CountDislikesByPostID(postID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"success":       true,
		"totalLikes":    likes,
		"totalDislikes": dislikes,
		"userReaction":  nil,
	}, nil
}

// GetUserReactionService récupère la réaction d'un utilisateur
func (s *postsService) GetUserReactionService(postID int64, userID int64) (map[string]interface{}, error) {
	value, err := s.postsRepo.GetUserReactionRepo(postID, userID)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"hasReacted": value != 0,
		"value":      nil,
	}

	if value != 0 {
		result["value"] = value
	}

	return result, nil
}

func (s *postsService) CreatePostService(post *domain.Post, categories []string) (int64, error) {
	if post.Title == "" || post.Content == "" {
		return 0, ErrInvalidInput
	}

	postID, err := s.postsRepo.CreatePostRepo(post)
	if err != nil {
		return 0, err
	}

	for _, catName := range categories {
		// Check if category exists
		catID, err := s.postsRepo.GetCategoryByNameRepo(catName)
		if err != nil {
			// Create category if not exists
			catID, err = s.postsRepo.CreateCategoryRepo(catName)
			if err != nil {
				return 0, err
			}
		}
		// Link category to post
		err = s.postsRepo.AddCategoryToPostRepo(postID, catID)
		if err != nil {
			return 0, err
		}
	}

	return postID, nil
}

func (s *postsService) GetAllCategoriesService() ([]string, error) {
	return s.postsRepo.GetAllCategoriesRepo()
}
