package server

import (
	"social-network/core/interfaces"
	"social-network/server/websocket"
)

type HttpAdapter struct {
	service interfaces.Core
	manager *websocket.Manager
}

func NewHttpAdapter(service interfaces.Core, manager *websocket.Manager) *HttpAdapter {
	return &HttpAdapter{
		service: service,
		manager: manager,
	}
}
