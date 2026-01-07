package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"real-time-forum/backend/domain"
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

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(NewUser)

	fmt.Println(NewUser)
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
