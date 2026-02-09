package repositories

import (
	"database/sql"
	"fmt"
	"real-time-forum/backend/domain"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) domain.MessageRepository {
	return &MessageRepo{db: db}
}

// CreateMessage insère un nouveau message dans la base de données
func (r *MessageRepo) CreateMessage(message *domain.Message) error {
	query := `
		INSERT INTO messages (sender_id, receiver_id, content, created_at, is_read)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, message.SenderID, message.ReceiverID, message.Content, message.CreatedAt, message.IsRead)
	if err != nil {
		return fmt.Errorf("error creating message: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting message ID: %w", err)
	}

	message.ID = id
	return nil
}

// GetMessagesBetweenUsers récupère les messages entre deux utilisateurs avec pagination
func (r *MessageRepo) GetMessagesBetweenUsers(userID1, userID2 int64, limit, offset int) ([]domain.Message, error) {
	query := `
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.created_at, m.is_read,
		       sender.username as sender_username, receiver.username as receiver_username
		FROM messages m
		JOIN users sender ON m.sender_id = sender.id
		JOIN users receiver ON m.receiver_id = receiver.id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, userID1, userID2, userID2, userID1, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying messages: %w", err)
	}
	defer rows.Close()

	var messages []domain.Message
	for rows.Next() {
		var msg domain.Message
		err := rows.Scan(
			&msg.ID,
			&msg.SenderID,
			&msg.ReceiverID,
			&msg.Content,
			&msg.CreatedAt,
			&msg.IsRead,
			&msg.SenderUsername,
			&msg.ReceiverUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning message: %w", err)
		}
		messages = append(messages, msg)
	}

	// Inverser l'ordre pour avoir les messages du plus ancien au plus récent
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, nil
}

// GetChatUsers récupère la liste des utilisateurs avec lesquels l'utilisateur a échangé des messages
func (r *MessageRepo) GetChatUsers(userID int64) ([]domain.ChatUser, error) {
	query := `
		WITH user_messages AS (
			SELECT 
				CASE 
					WHEN sender_id = ? THEN receiver_id 
					ELSE sender_id 
				END as other_user_id,
				content as last_message,
				created_at as last_message_at,
				CASE 
					WHEN receiver_id = ? AND is_read = 0 THEN 1 
					ELSE 0 
				END as is_unread
			FROM messages
			WHERE sender_id = ? OR receiver_id = ?
		),
		latest_messages AS (
			SELECT 
				other_user_id,
				last_message,
				last_message_at,
				SUM(is_unread) as unread_count,
				ROW_NUMBER() OVER (PARTITION BY other_user_id ORDER BY last_message_at DESC) as rn
			FROM user_messages
			GROUP BY other_user_id, last_message, last_message_at
		)
		SELECT 
			u.id,
			u.username,
			COALESCE(lm.last_message, '') as last_message,
			COALESCE(lm.last_message_at, '') as last_message_at,
			COALESCE(lm.unread_count, 0) as unread_count
		FROM users u
		LEFT JOIN latest_messages lm ON u.id = lm.other_user_id AND lm.rn = 1
		WHERE u.id != ? AND (lm.other_user_id IS NOT NULL OR u.id IN (
			SELECT DISTINCT sender_id FROM messages WHERE receiver_id = ?
			UNION
			SELECT DISTINCT receiver_id FROM messages WHERE sender_id = ?
		))
		ORDER BY lm.last_message_at DESC, u.username ASC
	`

	rows, err := r.db.Query(query, userID, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying chat users: %w", err)
	}
	defer rows.Close()

	var chatUsers []domain.ChatUser
	for rows.Next() {
		var user domain.ChatUser
		var lastMessageAt sql.NullString
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.LastMessage,
			&lastMessageAt,
			&user.UnreadCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat user: %w", err)
		}

		// Parse la date si elle existe
		if lastMessageAt.Valid && lastMessageAt.String != "" {
			// SQLite retourne les dates au format string
			// On pourrait parser ici si nécessaire
		}

		chatUsers = append(chatUsers, user)
	}

	return chatUsers, nil
}

// MarkMessagesAsRead marque tous les messages d'un expéditeur vers un destinataire comme lus
func (r *MessageRepo) MarkMessagesAsRead(senderID, receiverID int64) error {
	query := `
		UPDATE messages
		SET is_read = 1
		WHERE sender_id = ? AND receiver_id = ? AND is_read = 0
	`

	_, err := r.db.Exec(query, senderID, receiverID)
	if err != nil {
		return fmt.Errorf("error marking messages as read: %w", err)
	}

	return nil
}

// GetUnreadCount retourne le nombre de messages non lus pour un utilisateur
func (r *MessageRepo) GetUnreadCount(receiverID int64) (int, error) {
	query := `
		SELECT COUNT(*) FROM messages
		WHERE receiver_id = ? AND is_read = 0
	`

	var count int
	err := r.db.QueryRow(query, receiverID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting unread count: %w", err)
	}

	return count, nil
}
