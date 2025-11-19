package handlers

import (
	"net/http"

	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
)

// DashboardHandler handles dashboard endpoints
type DashboardHandler struct {
	dashboardService *service.DashboardService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		dashboardService: dashboardService,
	}
}

// GetStats returns dashboard statistics
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.dashboardService.GetStats(r.Context())
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve dashboard stats")
		return
	}

	httputil.RespondSuccess(w, http.StatusOK, stats)
}

// GetRecentActivity returns recent activity
func (h *DashboardHandler) GetRecentActivity(w http.ResponseWriter, r *http.Request) {
	activities, err := h.dashboardService.GetRecentActivity(r.Context(), 10)
	if err != nil {
		httputil.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve recent activity")
		return
	}

	httputil.RespondSuccess(w, http.StatusOK, activities)
}
