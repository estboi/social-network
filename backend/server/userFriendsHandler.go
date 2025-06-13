package server

import (
	"encoding/json"
	"log"
	"net/http"
	"social-network/core/entities"
	"social-network/server/websocket"
	"social-network/sessions"

	routing "github.com/Zewasik/go-router"
)

func (handler *HttpAdapter) ActionHandlerPostSubscribeOnUser(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	otherUser, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		return
	}

	// Step 3: Call the action handler and get the response and error.
	response, err := handler.service.SubscribeProcessing(userId, otherUser)
	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) and check for errors.
	encodedResponse, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		// Handle the encoding error and send an error response.
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}

	var Event websocket.Event
	Event.Type = "New_Notification"
	notification := entities.NotificationVM{UserId: otherUser, SourceId: userId, SourceType: "friend", NotifType: "friendRequest", Content: "You`ve got new friend request"}
	if err := handler.service.NotificationRecord(notification); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	encodedNotification, _ := json.Marshal(notification)
	Event.Payload = json.RawMessage(encodedNotification)
	for c := range handler.manager.Clients {
		if c.Id == otherUser {
			c.MessageChan <- Event
		}
	}

	// Step 5: Send the encoded response as an HTTP response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encodedResponse)
}

func (handler *HttpAdapter) ActionHandlerPostUnsubscribeOnUser(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	otherUser, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		log.Println(err)
		return
	}
	// Step 3: Call the action handler and get the response and error.
	response, err := handler.service.UnsubscribeProcessing(userId, otherUser)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) and check for errors.
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

func (handler *HttpAdapter) UsersHandlerModify(w http.ResponseWriter, r *http.Request) {
	userId, err := sessions.Validate(r)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response, err := handler.service.ModifyUserProcessing(userId)

	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: Encode the response entity (assuming it's a struct) and check for errors.
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
