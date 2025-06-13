package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/sessions"
)

func (handler *HttpAdapter) NotificationsHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Step 1: Skip decoding input entity 'stop' since it's not required in this case

	// Step 2: Check for decoding error (since there's no decoding, there are no errors to check)
	userId, _ := sessions.Validate(r)
	// Step 3: Call the handler.service.NotificationsHandlerGet
	response, err := handler.service.NotificationsHandlerGet(userId)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < NotificationVM and check for errors
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
