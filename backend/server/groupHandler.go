package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"social-network/core/entities"
	"social-network/server/websocket"
	"social-network/sessions"
	"strconv"

	routing "github.com/Zewasik/go-router"
)

// GET
func (handler *HttpAdapter) GroupsHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	// Step 1: get userID
	userID, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: processing
	groups, err := handler.service.GroupsAllProc(userID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: encode
	encodedResponse, err := json.Marshal(groups)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 4: send
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}

func (handler *HttpAdapter) GroupsHandlerGetConnected(w http.ResponseWriter, r *http.Request) {
	userID, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 3: Call the Groups handler and get the response and error.
	response, err := handler.service.GroupsHandlerGetConnected(userID)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) as is and check for errors.
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

func (handler *HttpAdapter) GroupsHandlerGetCreated(w http.ResponseWriter, r *http.Request) {
	// Step 1: get userID
	userID, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 3: Call the Groups handler and get the response and error.
	response, err := handler.service.GroupsHandlerGetCreated(userID)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) as is and check for errors.
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

func (handler *HttpAdapter) GroupsHandlerGetProfile(w http.ResponseWriter, r *http.Request) {
	UserId, _ := sessions.Validate(r)

	// Step 1: Get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Step 2: Call the handler.service.GroupsHandlerGetProfile
	response, err := handler.service.GroupsProfileProc(UserId, groupId)

	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusNotAcceptable)
		return
	}

	// Step 3: Simply encode the response entity < GroupVM and check for errors
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		// Handle the encoding error, for example, by sending an error response
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Step 4: Send the encoded response as an HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(encodedResponse)
	if err != nil {
		// Handle the write error, for example, by logging it
		fmt.Printf("Error writing response: %v\n", err)
	}
}

// GROUP OPTIONS
func (handler *HttpAdapter) GroupsHandlerGetNotMembers(w http.ResponseWriter, r *http.Request) {
	// Step 1: get groupID
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: processing
	users, err := handler.service.GroupsGetNotMembersProc(groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Step 3: Decode and send
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Step 4: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}

func (handler *HttpAdapter) GroupsHandlerGetRequested(w http.ResponseWriter, r *http.Request) {
	// Step 1: get userId
	userId, _ := sessions.Validate(r)

	// Step 2: get groupID
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	users, err := handler.service.GroupsGetRequestedProc(userId, groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
	}

	// Step 4: Decode and send
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	// Step 5: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}

// CREATION
func (handler *HttpAdapter) GroupsHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode input entity 'GroupDTO'
	var inputEntity entities.GroupDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&inputEntity); err != nil {
		// Handle the decoding error, for example, by sending an error response
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 2: Get the userId
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Step 3: Call the handler.service.GroupsHandlerPost
	response, err := handler.service.GroupsCreateProc(userId, inputEntity)
	if err != nil {
		// Handle the error, for example, by sending an error response
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

// ACTIONS
func (handler *HttpAdapter) GroupsRequestHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get userId
	userId, _ := sessions.Validate(r)

	// Step 2: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	CreatorId, err := handler.service.GroupsRequestProc(userId, groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	// Step 4: send notification about invitation
	var Event websocket.Event
	Event.Type = "New_Notification"

	notification := entities.NotificationVM{UserId: CreatorId, SourceId: groupId, SourceType: "groups", NotifType: "request", Content: "You have request to group"}
	if err := handler.service.NotificationRecord(notification); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	encodedNotification, _ := json.Marshal(notification)
	Event.Payload = json.RawMessage(encodedNotification)
	for c := range handler.manager.Clients {
		if c.Id == CreatorId {
			c.MessageChan <- Event
		}
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) GroupsInviteHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Step 2: get invited UserId
	var invitedUserId entities.GroupsAction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&invitedUserId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	if err := handler.service.GroupsInviteProc(invitedUserId.UserID, groupId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Step 4: send notification about invitation
	var Event websocket.Event
	Event.Type = "New_Notification"

	notification := entities.NotificationVM{UserId: invitedUserId.UserID, SourceId: groupId, SourceType: "groups", NotifType: "invite", Content: "You have invited to group"}
	if err := handler.service.NotificationRecord(notification); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	encodedNotification, _ := json.Marshal(notification)
	Event.Payload = json.RawMessage(encodedNotification)
	for c := range handler.manager.Clients {
		if c.Id == invitedUserId.UserID {
			c.MessageChan <- Event
		}
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) GroupsAcceptHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get RequestingUserId
	var RequestingUserId entities.GroupsAction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&RequestingUserId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	if err := handler.service.GroupsAcceptProc(RequestingUserId.UserID, groupId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) GroupsDenyHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get RequestingUserId
	var RequestingUserId entities.GroupsAction
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&RequestingUserId); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	if err := handler.service.GroupsDenyProc(RequestingUserId.UserID, groupId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) GroupsAcceptInviteHandler(w http.ResponseWriter, r *http.Request) {

	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get RequestingUserId
	userId, _ := sessions.Validate(r)

	// Step 3: processing
	if err := handler.service.GroupsAcceptProc(userId, groupId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) GroupsInviteDenyHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get RequestingUserId
	userId, _ := sessions.Validate(r)

	// Step 3: processing
	if err := handler.service.GroupsDenyProc(userId, groupId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	// Step 4 : send accept status
	w.WriteHeader(http.StatusAccepted)
}

// IMAGES HANDLER
func (handler HttpAdapter) ImageGroupCreateHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: unparse image
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

	// Step 3: call proc method
	if err := handler.service.GroupsImageCreateProc(groupId, content); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler HttpAdapter) ImageGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Step 1: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: call the Proc method
	content, err := handler.service.GroupsImageProc(groupId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Decode content
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
