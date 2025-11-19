package handlers

import (
"net/http"
"time"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/pkg/httputil"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

type WastageHandler struct {
wastageService *service.WastageService
}

func NewWastageHandler(wastageService *service.WastageService) *WastageHandler {
return &WastageHandler{
wastageService: wastageService,
}
}

func (h *WastageHandler) RecordWastage(w http.ResponseWriter, r *http.Request) {
var wastage domain.WastageLog
if err := httputil.DecodeJSON(r, &wastage); err != nil {
httputil.SendError(w, err)
return
}

if err := h.wastageService.RecordWastage(r.Context(), &wastage); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusCreated, wastage)
}

func (h *WastageHandler) GetByID(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

wastage, err := h.wastageService.GetWastage(r.Context(), id)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, wastage)
}

func (h *WastageHandler) List(w http.ResponseWriter, r *http.Request) {
startDateStr := r.URL.Query().Get("startDate")
endDateStr := r.URL.Query().Get("endDate")

var startDate, endDate time.Time
if startDateStr != "" {
startDate, _ = time.Parse("2006-01-02", startDateStr)
}
if endDateStr != "" {
endDate, _ = time.Parse("2006-01-02", endDateStr)
}

wastages, err := h.wastageService.ListWastage(r.Context(), startDate, endDate)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, wastages)
}

func (h *WastageHandler) Update(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var wastage domain.WastageLog
if err := httputil.DecodeJSON(r, &wastage); err != nil {
httputil.SendError(w, err)
return
}

wastage.ID = id

if err := h.wastageService.UpdateWastage(r.Context(), &wastage); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, wastage)
}

func (h *WastageHandler) Delete(w http.ResponseWriter, r *http.Request) {
id := mux.Vars(r)["id"]

var req domain.DeleteRequest
if err := httputil.DecodeJSON(r, &req); err != nil {
httputil.SendError(w, err)
return
}

if err := h.wastageService.DeleteWastage(r.Context(), id, &req); err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "wastage entry deleted successfully"})
}

func (h *WastageHandler) GetReport(w http.ResponseWriter, r *http.Request) {
startDateStr := r.URL.Query().Get("startDate")
endDateStr := r.URL.Query().Get("endDate")

var startDate, endDate time.Time
if startDateStr != "" {
startDate, _ = time.Parse("2006-01-02", startDateStr)
}
if endDateStr != "" {
endDate, _ = time.Parse("2006-01-02", endDateStr)
}

totalCost, err := h.wastageService.CalculateTotalWastageCost(r.Context(), startDate, endDate)
if err != nil {
httputil.SendError(w, err)
return
}

httputil.SendJSON(w, http.StatusOK, map[string]float64{"totalCost": totalCost})
}
