package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/gestionigrillo/secure-samsung-vault/internal/service"
)

type AdminHandler struct {
	commandSvc *service.CommandService
}

func NewAdminHandler(commandSvc *service.CommandService) *AdminHandler {
	return &AdminHandler{commandSvc: commandSvc}
}

func (h *AdminHandler) CreateCommand(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeviceID == "" || req.Type == "" {
		writeError(w, http.StatusBadRequest, "device_id and type required")
		return
	}

	// Validate command type
	switch req.Type {
	case models.CommandLock, models.CommandWipe, models.CommandRevoke:
		// valid
	default:
		writeError(w, http.StatusBadRequest, "invalid command type: use lock, wipe, or revoke")
		return
	}

	cmd, err := h.commandSvc.CreateCommand(r.Context(), &req, remoteIP(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, cmd)
}
