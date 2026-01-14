package services

import (
	"real-time-forum/backend/domain"
)

type commentsService struct {
	commentsRepo domain.CommentsRepository
}

func NewCommentsService(cr domain.CommentsRepository) domain.CommentsService {
	return &commentsService{commentsRepo: cr}
}

func (s *commentsService) GetCommentsByPostIDService(postID int64) ([]*domain.Comment, error) {
	comments, err := s.commentsRepo.GetCommentsByPostIDRepo(postID)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		username, err := s.commentsRepo.GetUsernameByIdRepo(comment.UserID)
		if err != nil {
			return nil, err
		}
		comment.Username = username
	}
	return comments, nil
}

func (s *commentsService) NewCommentService(comment *domain.Comment) error {
	err := s.commentsRepo.NewCommentRepo(comment)
	if err != nil {
		return err
	}
	return nil
}
