package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
)

// ForecastHandler handles demand forecasting endpoints
type ForecastHandler struct {
	forecastService *service.ForecastService
}

// NewForecastHandler creates a new forecast handler
func NewForecastHandler(forecastService *service.ForecastService) *ForecastHandler {
	return &ForecastHandler{
		forecastService: forecastService,
	}
}

// ForecastRequest represents a forecasting request
type ForecastRequest struct {
	ItemName   string `json:"itemName"`
	Days       int    `json:"days"`
	Historical bool   `json:"historical"`
}

// GenerateForecast handles AI-powered demand forecasting
func (h *ForecastHandler) GenerateForecast(w http.ResponseWriter, r *http.Request) {
	var req ForecastRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request format")
		return
	}

	if req.ItemName == "" {
		httputil.RespondError(w, http.StatusBadRequest, "MISSING_ITEM_NAME", "itemName is required")
		return
	}

	forecast, err := h.forecastService.GenerateForecast(r.Context(), service.ForecastRequest{
		ItemName:   req.ItemName,
		Days:       req.Days,
		Historical: req.Historical,
	})

	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "FORECAST_FAILED", "Failed to generate forecast")
		return
	}

	httputil.RespondSuccess(w, http.StatusOK, forecast)
}
