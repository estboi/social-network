package server

import (
	"encoding/json"
	"net/http"
	"social-network/core/entities"
	"social-network/server/websocket"
	"social-network/sessions"

	routing "github.com/Zewasik/go-router"
)

func (handler *HttpAdapter) EventHandlerGetAll(w http.ResponseWriter, r *http.Request) {
	// Step 1: Get the UserId
	UserId, _ := sessions.Validate(r)

	// Step 2: Call the Event handler and get the response and error.
	response, err := handler.service.EventGetAllProc(UserId)
	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Encode the response entity (assuming it's a struct) as is and check for errors.
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

func (handler *HttpAdapter) EventHandlerGetGroup(w http.ResponseWriter, r *http.Request) {
	// Step 1: get UserId
	userId, _ := sessions.Validate(r)

	// Step 2: get groupId
	groupId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// Step 3: Call the Event handler and get the response and error.
	response, err := handler.service.EventGetGroupProc(userId, groupId)
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

func (handler *HttpAdapter) EventHandlerPost(w http.ResponseWriter, r *http.Request) {
	// Step 1: Decode input entity "EventDTO".
	var eventDTO entities.EventDTO
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&eventDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 2: get UserID
	userId, _ := sessions.Validate(r)

	// Step 3: Call the Event handler and get the response and error.
	users, err := handler.service.EventCreateProc(userId, eventDTO)
	if err != nil {
		// Handle the error and send an error response.
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 4: send notification about invitation
	var Event websocket.Event
	Event.Type = "New_Notification"
	for _, v := range users {
		notification := entities.NotificationVM{UserId: v, SourceId: eventDTO.GroupId, SourceType: "event", NotifType: "newEvent", Content: "New event!"}
		if err := handler.service.NotificationRecord(notification); err != nil {
			http.Error(w, err.Error(), http.StatusNotAcceptable)
			return
		}

		encodedNotification, _ := json.Marshal(notification)
		Event.Payload = json.RawMessage(encodedNotification)
		for c := range handler.manager.Clients {
			if c.Id == v {
				c.MessageChan <- Event
			}
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func (handler *HttpAdapter) ActionHandlerPostJoinEvent(w http.ResponseWriter, r *http.Request) {
	// Step 1: get UserId
	userId, _ := sessions.Validate(r)

	// Step 2: get EventId
	eventId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	if err := handler.service.EventAcceptProc(userId, eventId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(202)
}

func (handler *HttpAdapter) ActionHandlerPostDenyEvent(w http.ResponseWriter, r *http.Request) {
	// Step 1: get UserId
	userId, _ := sessions.Validate(r)

	// Step 2: get EventId
	eventId, err := routing.GetRequestParamInt(r, "id")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Step 3: processing
	if err := handler.service.EventDenyProc(userId, eventId); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(202)
}
