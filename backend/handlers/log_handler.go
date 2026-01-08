package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/domain"
	"time"
)

var clientService domain.ClientService

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var NewUser domain.User
	err = json.NewDecoder(r.Body).Decode(&NewUser)
	fmt.Println(NewUser.Email, NewUser.Username, NewUser.Password)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, exist := clientService.Authentification(NewUser.Username, NewUser.Email, NewUser.Password)
	if !exist {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		fmt.Println(err)
		return
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		fmt.Println("❌ error generating token ")
	}

	sessionToken := base64.URLEncoding.EncodeToString(bytes)
	fmt.Println(sessionToken)

	// SetCookie AVANT WriteHeader (important !)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true,  // À activer quand vous passerez en HTTPS
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewUser)

	err = clientService.UpdateTokenService(NewUser.Username, NewUser.Email, sessionToken)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("connected user:", NewUser.Email, NewUser.Username)
	//http.Redirect(w, r, "/", http.StatusSeeOther)
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var NewUser domain.User
	err = json.NewDecoder(r.Body).Decode(&NewUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = clientService.Register(&NewUser)
	if err != nil {
		http.Error(w, "Error registering user", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(NewUser)
	fmt.Println("User registered:", NewUser)
}
