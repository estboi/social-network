package websocket

import (
	"encoding/json"
	"log"
	"social-network/core/entities"
	"time"
)

func SendMessage(event Event, c *Client) error {
	var data entities.Message
	if err := json.Unmarshal(event.Payload, &data); err != nil {
		log.Printf("error unmarshaling event: %s\n", err)
		return err
	}
	data.Time = time.Now().UTC().Format("02.01.2006")
	data.SenderId = c.Id

	err := c.Manager.service.SendMessage(data)
	if err != nil {
		log.Printf("Error Sending message to core: %s\n", err)
		return err
	}

	//Transform data to RawJSON
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("error marshaling data: %s\n", err)
		return err
	}
	sendData := json.RawMessage(dataBytes)
	event.Payload = sendData

	//Send message to chat members
	if err != nil {
		log.Printf("Error getting chat members: %s", err)
	}

	for client := range c.Manager.Clients {
		client.MessageChan <- event
	}
	return nil
}
