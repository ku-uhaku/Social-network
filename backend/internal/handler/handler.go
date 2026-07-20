package handler

import (
	"kuu/internal/service"
	"kuu/internal/websocket"
)

type Handler struct {
	Service *service.Service
	Hub     *websocket.Hub
}

func New(svc *service.Service, hub *websocket.Hub) *Handler {
	return &Handler{
		Service: svc,
		Hub:     hub,
	}
}
