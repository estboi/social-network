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

// GET HANDLER
func (handler *HttpAdapter) CommentsGetHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: Get PostID
	targetPostId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: call the processing function
	comments, err := handler.service.CommentsGetProc(targetPostId)
	if err != nil {
		log.Println(err)
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: encode response
	encodedResponse, err := json.Marshal(comments)
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

// POST HANDLER
func (handler *HttpAdapter) CommentsPostHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	// Step 1: Decode input entity 'PostDTO'
	var inputEntity entities.CommentDTO
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&inputEntity.Content); err != nil {
		log.Println(err)
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 2: Get the group ID
	inputEntity.PostId, err = routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 3: Get the user ID
	inputEntity.UserId, err = sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Step 3: Call the handler.service.PostsHandlerPost
	response, err := handler.service.CommentPostProc(inputEntity)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Simply encode the response entity < Status and check for errors
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

// IMAGES HANDLERS
func (handler *HttpAdapter) CommentsImageGetHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get target comment ID
	commentID, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: call the proc func
	response, err := handler.service.CommentsImageGetProc(commentID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Decode the image into JSON
	// Set the Content-Length header to specify the size of the file
	w.Header().Set("Content-Length", strconv.Itoa(len(response)))

	// Write the byte array as the response body
	if len(response) != 0 {
		_, err = w.Write(response)
		if err != nil {
			// Handle the write error, for example, by logging it
			fmt.Printf("error writing response: %v\n", err)
		}
	}
}
func (handler *HttpAdapter) CommentsImagePostHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get the commentId
	commentId, err := routing.GetRequestParamInt(r, "id")
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

	if err := handler.service.CommentsImageCreateProc(commentId, content); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

// ACTIONS HANDLER
func (handler *HttpAdapter) CommentsLikeHandler(w http.ResponseWriter, r *http.Request) {
	var votesData entities.VotesData
	var err error

	// Step 1: get userID
	votesData.UserId, err = sessions.Validate(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get groupID
	votesData.TargetId, err = routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: decode the g
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&votesData); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 4: record the like in database
	votesData, err = handler.service.CommentsLikeProc(votesData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 5: Send updated info back
	encodedResponse, err := json.Marshal(votesData)
	if err != nil {
		// Handle the encoding error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 6: Send the encoded response as an HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(encodedResponse)
	if err != nil {
		// Handle the write error, for example, by logging it
		fmt.Printf("Error writing response: %v\n", err)
	}
}
