package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
)

// ReportHandler handles report generation endpoints
type ReportHandler struct {
	reportService *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
	}
}

// ReportRequest represents a report generation request
type ReportRequest struct {
	Type       service.ReportType `json:"type"`
	StartDate  string             `json:"startDate,omitempty"`
	EndDate    string             `json:"endDate,omitempty"`
	CustomerID string             `json:"customerId,omitempty"`
}

// GenerateReport handles report generation based on type
func (h *ReportHandler) GenerateReport(w http.ResponseWriter, r *http.Request) {
	var req ReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.RespondError(w, http.StatusBadRequest, "INVALID_INPUT", "Invalid request format")
		return
	}

	ctx := r.Context()

	switch req.Type {
	case service.ReportTypeSales:
		report, err := h.reportService.GenerateSalesReport(ctx, req.StartDate, req.EndDate)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "REPORT_GENERATION_FAILED", "Failed to generate sales report")
			return
		}
		httputil.RespondSuccess(w, http.StatusOK, report)

	case service.ReportTypeInventory:
		report, err := h.reportService.GenerateInventoryReport(ctx)
		if err != nil {
			httputil.RespondError(w, http.StatusInternalServerError, "REPORT_GENERATION_FAILED", "Failed to generate inventory report")
			return
		}
		httputil.RespondSuccess(w, http.StatusOK, report)

	case service.ReportTypeCustomer:
		if req.CustomerID == "" {
			httputil.RespondError(w, http.StatusBadRequest, "MISSING_CUSTOMER_ID", "customerId is required for customer reports")
			return
		}
		report, err := h.reportService.GenerateCustomerReport(ctx, req.CustomerID, req.StartDate, req.EndDate)
		if err != nil {
			if err == domain.ErrNotFound {
				httputil.RespondError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", "Customer not found")
				return
			}
			httputil.RespondError(w, http.StatusInternalServerError, "REPORT_GENERATION_FAILED", "Failed to generate customer report")
			return
		}
		httputil.RespondSuccess(w, http.StatusOK, report)

	default:
		httputil.RespondError(w, http.StatusBadRequest, "UNSUPPORTED_REPORT_TYPE", "Unsupported report type")
	}
}
