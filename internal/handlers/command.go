package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/service"
)

type CommandHandler struct {
	commandSvc *service.CommandService
}

func NewCommandHandler(commandSvc *service.CommandService) *CommandHandler {
	return &CommandHandler{commandSvc: commandSvc}
}

func (h *CommandHandler) AcknowledgeCommand(w http.ResponseWriter, r *http.Request) {
	commandIDStr := chi.URLParam(r, "commandId")
	commandID, err := uuid.Parse(commandIDStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid command ID"})
		return
	}

	var req models.AckCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.commandSvc.AcknowledgeCommand(r.Context(), commandID, req, r.RemoteAddr)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to acknowledge command"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}
