package handlers

import (
	"net/http"

	"github.com/example/pgvms/internal/domain"
	"github.com/example/pgvms/internal/pkg/httputil"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

type CrateHandler struct {
	crateService *service.CrateService
}

func NewCrateHandler(crateService *service.CrateService) *CrateHandler {
	return &CrateHandler{
		crateService: crateService,
	}
}

func (h *CrateHandler) ListAll(w http.ResponseWriter, r *http.Request) {
	crates, err := h.crateService.GetAllCrates(r.Context())
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, crates)
}

func (h *CrateHandler) IssueCrates(w http.ResponseWriter, r *http.Request) {
	var crate domain.CrateEntry
	if err := httputil.DecodeJSON(r, &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.crateService.IssueCrates(r.Context(), &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, crate)
}

func (h *CrateHandler) ReturnCrates(w http.ResponseWriter, r *http.Request) {
	var crate domain.CrateEntry
	if err := httputil.DecodeJSON(r, &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.crateService.ReturnCrates(r.Context(), &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusCreated, crate)
}

func (h *CrateHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerId"]

	balance, err := h.crateService.GetCrateBalance(r.Context(), customerID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]int{"balance": balance})
}

func (h *CrateHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["customerId"]

	crates, err := h.crateService.GetCrateHistory(r.Context(), customerID)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, crates)
}

func (h *CrateHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	crate, err := h.crateService.GetCrate(r.Context(), id)
	if err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, crate)
}

func (h *CrateHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var crate domain.CrateEntry
	if err := httputil.DecodeJSON(r, &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	crate.ID = id

	if err := h.crateService.UpdateCrate(r.Context(), &crate); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, crate)
}

func (h *CrateHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var req domain.DeleteRequest
	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	if err := h.crateService.DeleteCrate(r.Context(), id, &req); err != nil {
		httputil.SendError(w, err)
		return
	}

	httputil.SendJSON(w, http.StatusOK, map[string]string{"message": "crate entry deleted successfully"})
}
