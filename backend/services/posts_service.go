package services

import (
	"real-time-forum/backend/domain"
)

type postsService struct {
	postsRepo domain.PostsRepository
}

func NewPostsService(pr domain.PostsRepository) domain.PostsService {
	return &postsService{postsRepo: pr}
}
func (s *postsService) GetAllPostsService() ([]*domain.Post, error) {
	posts, err := s.postsRepo.GetAllPostsRepo()
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
		commentsCount, err := s.postsRepo.CountCommentsByPostID(post.ID)
		if err != nil {
			return nil, err
		}
		post.Likes = likesCount
		post.Comments = commentsCount
		post.Username = username
	}

	return posts, nil
}

func (s *postsService) GetPostByIDService(postID int64) (*domain.Post, error) {
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

	post.Likes = likesCount
	post.Comments = commentsCount
	post.Username = username
	post.Dislikes = dislikesCount
	return post, nil
}
