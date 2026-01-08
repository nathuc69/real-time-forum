package services

import (
	"fmt"
	"real-time-forum/backend/domain"

	"golang.org/x/crypto/bcrypt"
)

type clientService struct {
	repo domain.ClientRepository
}

func NewAuthService(repo domain.ClientRepository) domain.ClientService {
	return &clientService{repo: repo}
}

func (s *clientService) Authentification(username, email, password string) (*domain.User, bool) {
	user, found := s.repo.ClientLog(username, email)
	if !found {
		return nil, false
	}
	// Ici, vous pouvez ajouter la vérification du mot de passe si nécessaire
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Println("❌ invalid password")
		return nil, false
	}
	return user, true
}

func (s *clientService) Register(user *domain.User) error {
	// Hacher le mot de passe avant de le stocker
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("❌ error hashing password: %w", err)
	}
	user.Password = string(hashedPassword)

	err = s.repo.CreateClient(user)
	if err != nil {
		return fmt.Errorf("❌ error creating user: %w", err)
	}
	return nil
}
func (s *clientService) UpdateTokenService(username, email, token string) error {
	err := s.repo.UpdateTokenRepo(username, email, token)
	if err != nil {
		return fmt.Errorf("❌ error updating token in service: %w", err)
	}
	return nil
}

func (s *clientService) CheckTokenService(token string) (*domain.User, bool) {
	user, found := s.repo.CheckTokenRepo(token)
	if !found {
		return nil, false
	}
	return user, true
}

func (s *clientService) DeleteTokenService(token string) error {
	err := s.repo.DeleteTokenRepo(token)
	if err != nil {
		return fmt.Errorf("❌ error deleting token in service: %w", err)
	}
	return nil
}
