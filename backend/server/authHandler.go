package server

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/core/entities"
	"social-network/sessions"
)

func (handler *HttpAdapter) LoginHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode input entity 'LoginDTO'
	var inputEntity entities.LoginDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputEntity); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 2: Call the Login Processing Method
	userId, err := handler.service.LoginProcess(inputEntity)
	if err != nil {
		// Handle the error, for example, by sending an error response
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Setup session
	if err := sessions.NewSession(userId, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, err := sessions.Validate(r)
	if err != nil || userID <= 0 {
		log.Println("Error setting client ID", err)
	}
}

func (handler *HttpAdapter) RegisterHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode input entity 'RegisterDTO'
	var inputEntity entities.RegisterDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputEntity); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 2: Call the Registration processing method
	userId, err := handler.service.RegisterProc(inputEntity)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Setup new session token
	if err := sessions.NewSession(userId, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Check if user session is up so redirect
func CheckIfAuth(w http.ResponseWriter, r *http.Request) {
	if _, err := sessions.Validate(r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	} else {
		w.WriteHeader(202)
	}
}
