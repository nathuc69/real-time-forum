package handlers

import (
	"net/http"
	"real-time-forum/backend/domain"
	"real-time-forum/backend/middleware"
)

func Router(cs domain.ClientService) http.Handler {
	clientService = cs

	mux := http.NewServeMux()

	mux.Handle("/api/login", middleware.CORS(http.HandlerFunc(LoginHandler)))
	// mux.HandleFunc("/api/register", RegisterHandler)
	mux.Handle("/api/thread", middleware.CORS(http.HandlerFunc(RegisterHandler)))
	// Routes:
	// mux.Handle("/thread", middleware.Handle(http.HandlerFunc(ThreadHandler)))
	// mux.HandleFunc("/login", AuthenticateHandler)

	return middleware.CORS(mux)
}
