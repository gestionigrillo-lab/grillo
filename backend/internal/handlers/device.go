package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gestionigrillo/secure-samsung-vault/internal/models"
	"github.com/gestionigrillo/secure-samsung-vault/internal/service"
)

type DeviceHandler struct {
	deviceSvc  *service.DeviceService
	commandSvc *service.CommandService
}

func NewDeviceHandler(deviceSvc *service.DeviceService, commandSvc *service.CommandService) *DeviceHandler {
	return &DeviceHandler{deviceSvc: deviceSvc, commandSvc: commandSvc}
}

func (h *DeviceHandler) Enroll(w http.ResponseWriter, r *http.Request) {
	var req models.EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeviceModel == "" || req.AndroidVersion == "" || req.AppVersion == "" {
		writeError(w, http.StatusBadRequest, "device_model, android_version, app_version required")
		return
	}

	resp, err := h.deviceSvc.Enroll(r.Context(), &req, remoteIP(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *DeviceHandler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	var req models.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Use authenticated device ID
	if authID := r.Header.Get("X-Device-ID"); authID != "" {
		req.DeviceID = authID
	}

	if req.DeviceID == "" {
		writeError(w, http.StatusBadRequest, "device_id required")
		return
	}

	resp, err := h.deviceSvc.Heartbeat(r.Context(), &req, remoteIP(r))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *DeviceHandler) GetCommands(w http.ResponseWriter, r *http.Request) {
	deviceID := chi.URLParam(r, "deviceId")
	if deviceID == "" {
		writeError(w, http.StatusBadRequest, "device_id required")
		return
	}

	commands, err := h.commandSvc.GetPendingCommands(r.Context(), deviceID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if commands == nil {
		commands = []models.Command{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"commands": commands})
}

func (h *DeviceHandler) AckCommand(w http.ResponseWriter, r *http.Request) {
	commandID := chi.URLParam(r, "commandId")
	if commandID == "" {
		writeError(w, http.StatusBadRequest, "command_id required")
		return
	}

	var req models.AckCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.commandSvc.AcknowledgeCommand(r.Context(), commandID, &req, remoteIP(r)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
