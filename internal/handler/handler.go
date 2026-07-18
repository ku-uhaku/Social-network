package handler

import (
	"net/http"

	"kuu/internal/repository"
	"kuu/internal/websocket"
)

type Handler struct {
	Repo *repository.Repository
	Hub  *websocket.Hub
}

func New(repo *repository.Repository, hub *websocket.Hub) *Handler {
	return &Handler{
		Repo: repo,
		Hub:  hub,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
