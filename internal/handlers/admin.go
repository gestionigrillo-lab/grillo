package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/secure-samsung-vault/grillo/internal/middleware"
	"github.com/secure-samsung-vault/grillo/internal/models"
	"github.com/secure-samsung-vault/grillo/internal/service"
)

type AdminHandler struct {
	commandSvc *service.CommandService
}

func NewAdminHandler(commandSvc *service.CommandService) *AdminHandler {
	return &AdminHandler{commandSvc: commandSvc}
}

func (h *AdminHandler) IssueCommand(w http.ResponseWriter, r *http.Request) {
	var req models.IssueCommandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	issuedBy, _ := r.Context().Value(middleware.AdminUserKey).(string)

	resp, err := h.commandSvc.IssueCommand(r.Context(), req, issuedBy, r.RemoteAddr)
	if err != nil {
		if err.Error() == "deviceId is required" || contains(err.Error(), "invalid command type") {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if contains(err.Error(), "device not found") {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to issue command"})
		return
	}

	writeJSON(w, http.StatusCreated, resp)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
