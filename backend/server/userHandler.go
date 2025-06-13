package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/sessions"

	routing "github.com/Zewasik/go-router"
)

func (handler *HttpAdapter) UsersHandlerGetAll(w http.ResponseWriter, r *http.Request) {

	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 3: Call the handler.service.UsersHandlerGetAll
	response, err := handler.service.GetAllUsersProcessing(userId, from, to)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < UsersVM and check for errors
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

func (handler *HttpAdapter) UsersHandlerGetProfile(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	profileID, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		return
	}

	// Step 3: Call the handler.service.UsersHandlerGetProfile
	response, err := handler.service.GetUserProfileProcessing(userId, profileID)
	if err != nil {
		if err.Error() == "no access" {
			http.Error(w, "NO ACCESS", http.StatusNotAcceptable)
			return
		}
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if userId != profileID {
		response.CanModify = false
	} else {
		response.CanModify = true
	}

	// Step 4: Simply encode the response entity < UserFullVM and check for errors
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

func (handler *HttpAdapter) UsersHandlerGetFollowed(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Step 3: Call the handler.service.UsersHandlerGetFollowed
	response, err := handler.service.GetFollowedUsersProcessing(userId)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < UsersVM and check for errors
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

func (handler *HttpAdapter) UsersHandlerGetFollowers(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Step 3: Call the handler.service.UsersHandlerGetFollowers
	response, err := handler.service.GetFollowersProcessing(userId, from, to)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < UsersVM and check for errors
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
