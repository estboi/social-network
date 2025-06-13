package websocket

import (
	"errors"
	"log"
	"net/http"
	"social-network/core"
	"social-network/sessions"
	"sync"

	"github.com/gorilla/websocket"
)

const (
	NewMessage = "New_Message"
)

type Manager struct {
	service      core.Core
	Clients      ClientList // Connected users to WS
	sync.RWMutex            // Async blocking for write to map
	handlers     map[string]EventHandler
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func NewManager(s core.Core) *Manager {
	m := &Manager{
		Clients:  make(map[*Client]bool),
		handlers: make(map[string]EventHandler),
		service:  s,
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers["New_Message"] = SendMessage
}

func (m *Manager) routeEvent(event Event, c *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}
	return errors.New("there is no such event type")
}

func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	clientId, err := sessions.Validate(r)
	if err != nil && clientId <= 0 {
		// Handle the case when the ID is not available or invalid
		log.Println("Invalid or missing client ID. No websockets for U")
		return
	}
	log.Println("New connection")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection: %v\n", err)
		return
	}

	client := NewClient(ws, m)
	client.Id = clientId

	m.AddClient(client)

	go client.ReadLoop()
	go client.WriteLoop()

	close := ws.CloseHandler()

	ws.SetCloseHandler(func(code int, text string) error {
		close(code, text)
		return nil
	})
}

func (m *Manager) AddClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.Clients[client] = true
}

func (m *Manager) RemoveClient(Client *Client) {
	m.Lock()
	defer m.Unlock()
	Client.Connection.Close()
	delete(m.Clients, Client)
}

func checkOrigin(r *http.Request) bool {
	// origin := r.Header.Get("Origin")
	// if origin == "http://localhost:3000" || origin == "http://localhost:8080" {
	return true
	// }
	// return false
}
