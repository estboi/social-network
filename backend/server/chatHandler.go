package server

import (
	"encoding/json"
	"net/http"
	"social-network/sessions"

	routing "github.com/Zewasik/go-router"
)

func (handler *HttpAdapter) ChatsHandlerGetOneChat(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	otherUser, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		return
	}

	response, err := handler.service.GetChatMessages(userId, otherUser)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) as ChatVM and check for errors.

	encodedResponse, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		// Handle the encoding error and send an error response.
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}

	// Step 5: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}

func (handler *HttpAdapter) ChatsHandlerGetGroup(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	groupID, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		return
	}
	response, err := handler.service.ChatGetGroupMessagesProc(userId, groupID)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) as ChatVM and check for errors.
	encodedResponse, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		// Handle the encoding error and send an error response.
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}

	// Step 5: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}

func (handler *HttpAdapter) ChatsHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := handler.service.ChatGetAllMessagesProc(userId)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) as ChatVM and check for errors.
	encodedResponse, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		// Handle the encoding error and send an error response.
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}

	// Step 5: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}
