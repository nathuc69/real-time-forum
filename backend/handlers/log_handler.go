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

func SetClientService(cs domain.ClientService) {
	clientService = cs
}

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

func CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Vérifier le cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "No session found",
		})
		return
	}

	// Valider le token
	user, err := clientService.CheckTokenService(cookie.Value)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated": false,
			"message":       "Invalid or expired session",
		})
		return
	}

	// Session valide
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated": true,
		"username":      user.Username,
		"email":         user.Email,
	})
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Récupérer le cookie de session
	cookie, err := r.Cookie("session_token")
	if err != nil || cookie.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "No session found",
		})
		return
	}

	// Supprimer le token de la base de données
	err = clientService.DeleteTokenService(cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Error during logout",
		})
		fmt.Println(err)
		return
	}

	// Effacer le cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour), // Cookie expiré
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Logout successful",
	})

	fmt.Println("User logged out")
}
