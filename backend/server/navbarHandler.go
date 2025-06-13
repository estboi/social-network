package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/sessions"
)

var USERID, from, to = 0, 0, 10

func (handler *HttpAdapter) NavbarHandlerGet(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Step 3: Call the handler.service.NavbarHandlerGet
	response, err := handler.service.NavbarProc(userId)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < NavbarVM and check for errors
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		// Handle the encoding error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 5: Send the encoded response as an HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(encodedResponse)
	if err != nil {
		// Handle the write error, for example, by logging it
		fmt.Printf("Error writing response: %v\n", err)
	}
}
