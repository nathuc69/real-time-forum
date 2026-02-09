package services

import (
	"fmt"
	"real-time-forum/backend/domain"
	"time"
)

type MessageServ struct {
	repo domain.MessageRepository
}

func NewMessageService(repo domain.MessageRepository) domain.MessageService {
	return &MessageServ{repo: repo}
}

// SendMessage envoie un nouveau message
func (s *MessageServ) SendMessage(message *domain.Message) error {
	if message.Content == "" {
		return fmt.Errorf("message content cannot be empty")
	}

	if message.SenderID == 0 || message.ReceiverID == 0 {
		return fmt.Errorf("sender and receiver IDs are required")
	}

	if message.SenderID == message.ReceiverID {
		return fmt.Errorf("cannot send message to yourself")
	}

	// Définir la date de création si elle n'est pas définie
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now()
	}

	return s.repo.CreateMessage(message)
}

// GetConversation récupère la conversation entre deux utilisateurs
func (s *MessageServ) GetConversation(userID1, userID2 int64, limit, offset int) ([]domain.Message, error) {
	if userID1 == 0 || userID2 == 0 {
		return nil, fmt.Errorf("user IDs are required")
	}

	if limit <= 0 {
		limit = 10 // Limite par défaut
	}

	if limit > 50 {
		limit = 50 // Limite maximale
	}

	return s.repo.GetMessagesBetweenUsers(userID1, userID2, limit, offset)
}

// GetChatList récupère la liste des utilisateurs avec lesquels l'utilisateur a échangé
func (s *MessageServ) GetChatList(userID int64) ([]domain.ChatUser, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user ID is required")
	}

	return s.repo.GetChatUsers(userID)
}

// MarkAsRead marque les messages comme lus
func (s *MessageServ) MarkAsRead(senderID, receiverID int64) error {
	if senderID == 0 || receiverID == 0 {
		return fmt.Errorf("sender and receiver IDs are required")
	}

	return s.repo.MarkMessagesAsRead(senderID, receiverID)
}
