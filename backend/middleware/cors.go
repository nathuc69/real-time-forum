package middleware

import (
	"context"
	"net/http"
	"real-time-forum/backend/domain"
)

// allowOrigins whitelists frontend origins that may send credentials.
var allowOrigins = map[string]bool{
	"http://localhost:8000": true,
	"http://0.0.0.0:8000":   true,
	"http://127.0.0.1:3000": true,
}

var clientService domain.ClientService

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if allowOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetClientService(cs domain.ClientService) {
	clientService = cs
}

func HandleAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if clientService == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Vérifier le cookie de session à chaque requête
		cookie, err := r.Cookie("session_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "Unauthorized - No session", http.StatusUnauthorized)
			return
		}

		// Valider le token
		client, err := clientService.CheckTokenService(cookie.Value)
		if err != nil || client == nil {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		client.IsLogged = true

		// Stocker l'utilisateur dans le contexte
		ctx := context.WithValue(r.Context(), "user", client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
