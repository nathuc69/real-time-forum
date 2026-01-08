package domain

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Error     string `json:"error"`
}

type ClientRepository interface {
	ClientLog(username, email string) (*User, bool)
	CreateClient(user *User) error
	UpdateTokenRepo(username, email, token string) error
	CheckTokenRepo(token string) (*User, bool)
	DeleteTokenRepo(token string) error
}

type ClientService interface {
	Authentification(username, email, password string) (*User, bool)
	Register(user *User) error
	UpdateTokenService(username, email, token string) error
	CheckTokenService(token string) (*User, bool)
	DeleteTokenService(token string) error
}
