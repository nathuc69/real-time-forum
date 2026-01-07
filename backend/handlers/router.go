package handlers

import (
	"net/http"
	"real-time-forum/backend/domain"
)

func Router(cs domain.ClientService) http.Handler {
	clientService = cs

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", LoginHandler)
	mux.HandleFunc("/api/register", RegisterHandler)
	// Routes:
	// mux.Handle("/thread", middleware.Handle(http.HandlerFunc(ThreadHandler)))
	// mux.HandleFunc("/login", AuthenticateHandler)

	return mux
}
