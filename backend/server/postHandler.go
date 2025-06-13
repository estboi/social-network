package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/core/entities"
	"social-network/sessions"
	"strconv"

	routing "github.com/Zewasik/go-router"
)

func (handler *HttpAdapter) PostsHandlerGet(w http.ResponseWriter, r *http.Request) {
	// Step 1: Get postId
	postId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: Processing
	post, err := handler.service.PostGetProc(postId)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Encode
	encodedResponse, err := json.Marshal(post)
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

func (handler *HttpAdapter) PostsHandlerGetGroup(w http.ResponseWriter, r *http.Request) {
	// Step 1: Skip decoding input entity 'stop' since it's not required in this case
	currentUserId, _ := sessions.Validate(r)
	// Step 2: Check for decoding error (since there's no decoding, there are no errors to check)
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: Call the handler.service.PostsGetGroupProc
	response, err := handler.service.PostsGetGroupProc(groupId, currentUserId)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < PostVM and check for errors
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

func (handler *HttpAdapter) PostsHandlerGetHome(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "not unathorized user", http.StatusUnauthorized)
		return
	}

	// Step 3: Call the handler.service.PostsHandlerGetHome
	response, err := handler.service.PostsGetHomeProc(userId)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < PostVM and check for errors
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

func (handler *HttpAdapter) PostsHandlerGetUser(w http.ResponseWriter, r *http.Request) {
	currentUserId, _ := sessions.Validate(r)

	userId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := handler.service.PostsHandlerGetUser(userId, currentUserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	encodePosts, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(encodePosts)
}

// Handler for creating posts
func (handler *HttpAdapter) PostsHandlerPostPost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate
	userId, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, "not unathorized user", http.StatusUnauthorized)
		return
	}
	// Step 2: Check if its related to group
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		http.Error(w, "something went wrong", http.StatusBadRequest)
		return
	}

	// Step 3: Decode input entity 'PostDTO'
	var inputEntity entities.PostDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputEntity); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// Step 4: Call the handler.service.PostsHandlerPost
	postID, err := handler.service.PostPostProc(userId, groupId, inputEntity)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 4: encode response
	encodedResponse, err := json.Marshal(postID)
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

func (handler *HttpAdapter) PostActionVoteHandler(w http.ResponseWriter, r *http.Request) {
	var votesData entities.VotesData
	var err error

	// Step 1: get userId
	votesData.UserId, _ = sessions.Validate(r)

	// Step 1: Get postId that has been liked
	votesData.TargetId, err = routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 2: get votes number
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&votesData); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 3: Call the action proc
	votesData, err = handler.service.PostsLikeProc(votesData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 6: Send updated info back
	encodedResponse, err := json.Marshal(votesData)
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

// IMAGES HANDLER
func (handler *HttpAdapter) ImagePostCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Step 2: get the postId
	postId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	if err := handler.service.PostsImagePostProc(postId, content); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func (handler *HttpAdapter) ImagePostHandler(w http.ResponseWriter, r *http.Request) {
	// Step 2: get current post ID
	postId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call your core service method to get the byte array
	content, err := handler.service.PostsImageGetProc(postId)
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
