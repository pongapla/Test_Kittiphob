package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"test-kittiphob/backend/internal/repository"
)

type QueueHandler struct {
	repo *repository.Repository
}

func NewQueueHandler(repo *repository.Repository) *QueueHandler {
	return &QueueHandler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func (h *QueueHandler) Issue(w http.ResponseWriter, r *http.Request) {
	ticket, err := h.repo.Issue(r.Context())
	if errors.Is(err, repository.ErrQueueFull) {
		writeError(w, http.StatusConflict, "queue is full")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to issue ticket")
		return
	}
	writeJSON(w, http.StatusCreated, ticket)
}

func (h *QueueHandler) Current(w http.ResponseWriter, r *http.Request) {
	number, err := h.repo.Current(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read queue")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"queueNumber": number,
	})
}

func (h *QueueHandler) GetTicket(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ticket, err := h.repo.GetTicket(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "ticket not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read ticket")
		return
	}
	writeJSON(w, http.StatusOK, ticket)
}

func (h *QueueHandler) Reset(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.Reset(r.Context()); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to reset queue")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"queueNumber": "00",
		"message": "queue reset successfully",
	})
}
