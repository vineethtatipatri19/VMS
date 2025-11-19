package handlers

import (
	"net/http"
	"strconv"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

type ExpiryHandler struct {
	expiryService *service.ExpiryService
}

func NewExpiryHandler(expiryService *service.ExpiryService) *ExpiryHandler {
	return &ExpiryHandler{
		expiryService: expiryService,
	}
}

func (h *ExpiryHandler) GenerateAlerts(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DaysThreshold int `json:"DaysThreshold"`
	}

	// Try to decode JSON body first
	if err := httputil.DecodeJSON(r, &req); err == nil && req.DaysThreshold > 0 {
		// Use days from request body
	} else {
		// Fall back to query parameter
		daysStr := r.URL.Query().Get("days")
		req.DaysThreshold = 7 // default
		if daysStr != "" {
			if d, err := strconv.Atoi(daysStr); err == nil {
				req.DaysThreshold = d
			}
		}
	}

	if err := h.expiryService.GenerateAlerts(r.Context(), req.DaysThreshold); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, map[string]string{"message": "alerts generated successfully"})
}

func (h *ExpiryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	alert, err := h.expiryService.GetAlert(r.Context(), id)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, alert)
}

func (h *ExpiryHandler) List(w http.ResponseWriter, r *http.Request) {
	acknowledgedStr := r.URL.Query().Get("acknowledged")
	acknowledged := false
	if acknowledgedStr == "true" {
		acknowledged = true
	}

	alerts, err := h.expiryService.ListAlerts(r.Context(), acknowledged)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, alerts)
}

func (h *ExpiryHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	alerts, err := h.expiryService.GetPendingAlerts(r.Context())
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	// Wrap in object for consistent API response
	response := map[string]interface{}{
		"alerts": alerts,
	}
	httputil.SendJSON(w, http.StatusOK, response)
}

func (h *ExpiryHandler) Acknowledge(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req struct {
		AcknowledgedBy string `json:"acknowledgedBy"`
	}
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.expiryService.AcknowledgeAlert(r.Context(), id, req.AcknowledgedBy); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "alert acknowledged successfully"})
}

func (h *ExpiryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var alert domain.ExpiryAlert
	if err := httputil.DecodeJSON(r, &alert); err != nil {
		httputil.SendError(w, err)
		return
	}

	alert.ID = id

	if err := h.expiryService.UpdateAlert(r.Context(), &alert); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, alert)
}

func (h *ExpiryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req domain.DeleteRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.expiryService.DeleteAlert(r.Context(), id, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "alert deleted successfully"})
}
