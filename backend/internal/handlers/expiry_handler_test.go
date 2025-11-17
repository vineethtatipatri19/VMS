package handlers

import (
"net/http"
"net/http/httptest"
"testing"

"github.com/example/pgvms/internal/domain"
"github.com/example/pgvms/internal/service"
"github.com/gorilla/mux"
)

func TestExpiryHandler_GenerateAlerts(t *testing.T) {
mockRepo := &mockExpiryRepo{}
mockInvRepo := &mockInventoryRepo{}
expiryService := service.NewExpiryService(mockRepo, mockInvRepo)
handler := NewExpiryHandler(expiryService)

req := httptest.NewRequest(http.MethodPost, "/api/expiry/generate", nil)
w := httptest.NewRecorder()

handler.GenerateAlerts(w, req)

if w.Code != http.StatusCreated {
t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
}
}

func TestExpiryHandler_GetByID(t *testing.T) {
mockRepo := &mockExpiryRepo{
alert: &domain.ExpiryAlert{
ID:              "e1",
InventoryItemID: "i1",
},
}
mockInvRepo := &mockInventoryRepo{}
expiryService := service.NewExpiryService(mockRepo, mockInvRepo)
handler := NewExpiryHandler(expiryService)

req := httptest.NewRequest(http.MethodGet, "/api/expiry/e1", nil)
req = mux.SetURLVars(req, map[string]string{"id": "e1"})
w := httptest.NewRecorder()

handler.GetByID(w, req)

if w.Code != http.StatusOK {
t.Errorf("Expected 200, got %d", w.Code)
}
}
