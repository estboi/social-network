package server

import (
	"fmt"
	"io"
	"net/http"
	"social-network/sessions"
	"strconv"

	routing "github.com/Zewasik/go-router"
)

// USERS
func (handler *HttpAdapter) ImageUserCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get userId from session token
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Parse the form data to get the uploaded file
	err = r.ParseMultipartForm(10 << 20) // Max memory to use for parsing form (10 MB in this case)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file from the form data
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Read the file content into a []byte
	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := handler.service.ImageUserCreateProc(userId, content); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (handler *HttpAdapter) ImageUserHandler(w http.ResponseWriter, r *http.Request) {
	currentUserID, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "session is down", http.StatusUnauthorized)
		return
	}

	userId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no user present
	if userId == 0 {
		userId = currentUserID
	}

	// Call your core service method to get the byte array
	content, err := handler.service.ImageUserProc(userId)
	if err != nil {
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	// Set the Content-Length header to specify the size of the file
	w.Header().Set("Content-Length", strconv.Itoa(len(content)))

	// Write the byte array as the response body
	if len(content) != 0 {
		_, err = w.Write(content)
		if err != nil {
			// Handle the write error, for example, by logging it
			fmt.Printf("error writing response: %v\n", err)
		}
	}
}
