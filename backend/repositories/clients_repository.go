package repositories

import (
	"database/sql"
	"fmt"
	"real-time-forum/backend/domain"
)

type LogRepo struct {
	db *sql.DB
}

func NewLogRepository(db *sql.DB) domain.ClientRepository {
	return &LogRepo{db: db}
}

func (r *LogRepo) ClientLog(username, email string) (*domain.User, bool) {
	user := &domain.User{}

	err := r.db.QueryRow(`
		SELECT userName, email, password
		FROM users 
		WHERE username = ?`, username).Scan(&user.Username, &user.Email, &user.Password)
	err2 := r.db.QueryRow(`
		SELECT userName, email, password
		FROM users 
		WHERE email = ?`, email).Scan(&user.Username, &user.Email, &user.Password)
	if err == sql.ErrNoRows && err2 == sql.ErrNoRows {
		fmt.Println("❌ no user found with given username or email")
		return nil, false
		// Aucun utilisateur trouvé
	} else if err != nil && err2 != nil {
		// Erreur réelle (connexion, SQL, etc.)
		fmt.Println("Erreur SQL:", err)
		return nil, false
	}

	// Si on arrive ici, un utilisateur existe
	return user, true
}

func (r *LogRepo) UpdateTokenRepo(username, email, token string) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET token = ? 
		WHERE username = ? OR email = ?`, token, username, email)
	if err != nil {
		return fmt.Errorf("❌ error updating token: %w", err)
	}
	return nil
}

func (r *LogRepo) CheckTokenRepo(token string) (*domain.User, bool) {
	user := &domain.User{}

	err := r.db.QueryRow(`SELECT userName, email, id FROM users WHERE token = ?`, token).Scan(&user.Username, &user.Email, &user.ID)
	if err == sql.ErrNoRows {
		fmt.Println("❌ no user found with given token")
		return nil, false
	} else if err != nil {
		fmt.Println("Erreur SQL:", err)
		return nil, false
	}

	return user, true
}

func (r *LogRepo) DeleteTokenRepo(token string) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET token = NULL 
		WHERE token = ?`, token)
	if err != nil {
		return fmt.Errorf("❌ error deleting token: %w", err)
	}
	return nil
}

func (r *LogRepo) CreateClient(user *domain.User) error {
	_, err := r.db.Exec(`
		INSERT INTO users (username, password, email, age, firstName, lastName, gender)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.Username, user.Password, user.Email, user.Age, user.FirstName, user.LastName, user.Gender)
	if err != nil {
		return fmt.Errorf("❌ error inserting new user: %w", err)
	}
	return nil
}
