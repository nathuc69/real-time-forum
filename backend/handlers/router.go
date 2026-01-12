package handlers

import (
	"net/http"
	"real-time-forum/backend/domain"
	"real-time-forum/backend/middleware"
)

func Router(cs domain.ClientService) http.Handler {
	clientService = cs
	middleware.SetClientService(cs)

	mux := http.NewServeMux()

	// Routes publiques (pas d'authentification requise)
	mux.Handle("/api/login", middleware.CORS(http.HandlerFunc(LoginHandler)))
	mux.Handle("/api/register", middleware.CORS(http.HandlerFunc(RegisterHandler)))
	mux.Handle("/api/check-auth", middleware.CORS(http.HandlerFunc(CheckAuthHandler)))

	// Routes:
	// mux.Handle("/thread", middleware.Handle(http.HandlerFunc(ThreadHandler)))
	// mux.HandleFunc("/login", AuthenticateHandler)

	return middleware.CORS(mux)
}
