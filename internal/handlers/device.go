package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/service"
)

type DeviceHandler struct {
	deviceSvc  *service.DeviceService
	commandSvc *service.CommandService
}

func NewDeviceHandler(deviceSvc *service.DeviceService, commandSvc *service.CommandService) *DeviceHandler {
	return &DeviceHandler{
		deviceSvc:  deviceSvc,
		commandSvc: commandSvc,
	}
}

func (h *DeviceHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	var req models.EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.deviceSvc.Enroll(r.Context(), req, r.RemoteAddr)
	if err != nil {
		if err.Error() == "device already enrolled" {
			writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
			return
		}
		if err.Error() == "deviceId is required" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "enrollment failed"})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *DeviceHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req models.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	resp, err := h.deviceSvc.Heartbeat(r.Context(), req, r.RemoteAddr)
	if err != nil {
		if err.Error() == "deviceId is required" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "heartbeat failed"})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DeviceHandler) GetPendingCommands(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")
	if deviceID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "deviceId is required"})
		return
	}

	commands, err := h.commandSvc.GetPendingCommands(r.Context(), deviceID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get commands"})
		return
	}

	if commands == nil {
		commands = []models.DeviceCommand{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"deviceId": deviceID,
		"commands": commands,
	})
}
