package domain

import "time"

// Message représente un message privé entre deux utilisateurs
type Message struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"senderId"`
	ReceiverID int64     `json:"receiverId"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt"`
	IsRead     bool      `json:"isRead"`
	// Informations supplémentaires pour le frontend
	SenderUsername   string `json:"senderUsername,omitempty"`
	ReceiverUsername string `json:"receiverUsername,omitempty"`
}

// ChatUser représente un utilisateur dans la liste de chat
type ChatUser struct {
	ID             int64     `json:"id"`
	Username       string    `json:"username"`
	IsOnline       bool      `json:"isOnline"`
	LastMessage    string    `json:"lastMessage,omitempty"`
	LastMessageAt  time.Time `json:"lastMessageAt,omitempty"`
	UnreadCount    int       `json:"unreadCount"`
}

// MessageRepository définit les méthodes pour interagir avec la base de données
type MessageRepository interface {
	CreateMessage(message *Message) error
	GetMessagesBetweenUsers(userID1, userID2 int64, limit, offset int) ([]Message, error)
	GetChatUsers(userID int64) ([]ChatUser, error)
	MarkMessagesAsRead(senderID, receiverID int64) error
	GetUnreadCount(receiverID int64) (int, error)
}

// MessageService définit la logique métier pour les messages
type MessageService interface {
	SendMessage(message *Message) error
	GetConversation(userID1, userID2 int64, limit, offset int) ([]Message, error)
	GetChatList(userID int64) ([]ChatUser, error)
	MarkAsRead(senderID, receiverID int64) error
}
