package testutil

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/pgvms/internal/handlers"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/internal/router"
	"github.com/example/pgvms/internal/service"
	"github.com/gorilla/mux"
)

// E2ETestContext holds the complete test environment for E2E tests
type E2ETestContext struct {
	DB     *sql.DB
	Router *mux.Router
	Server *httptest.Server
	Client *http.Client
}

// SetupE2ETest creates a complete test environment with database, router, and handlers
func SetupE2ETest(t *testing.T) *E2ETestContext {
	t.Helper()

	// Setup test database
	db := SetupTestDB(t)

	// Initialize repositories
	customerRepo := postgres.NewCustomerRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	saleItemRepo := postgres.NewSaleItemRepository(db)
	crateRepo := postgres.NewCrateRepository(db)
	wastageRepo := postgres.NewWastageRepository(db)
	expiryRepo := postgres.NewExpiryAlertRepository(db)
	paymentRepo := postgres.NewPaymentScheduleRepository(db)

	// Initialize services
	customerService := service.NewCustomerService(customerRepo)
	inventoryService := service.NewInventoryService(inventoryRepo)
	transactionService := service.NewTransactionService(
		transactionRepo,
		customerRepo,
		inventoryRepo,
		saleItemRepo,
	)
	saleItemService := service.NewSaleItemService(saleItemRepo, inventoryRepo)
	crateService := service.NewCrateService(crateRepo, customerRepo)
	wastageService := service.NewWastageService(wastageRepo, inventoryRepo)
	expiryService := service.NewExpiryService(expiryRepo, inventoryRepo)
	paymentService := service.NewPaymentService(paymentRepo, transactionRepo, customerRepo)

	// Initialize handlers
	customerHandler := handlers.NewCustomerHandler(customerService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	saleItemHandler := handlers.NewSaleItemHandler(saleItemService)
	crateHandler := handlers.NewCrateHandler(crateService)
	wastageHandler := handlers.NewWastageHandler(wastageService)
	expiryHandler := handlers.NewExpiryHandler(expiryService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Create handlers struct
	routerHandlers := &router.Handlers{
		Customer:    customerHandler,
		Inventory:   inventoryHandler,
		Transaction: transactionHandler,
		SaleItem:    saleItemHandler,
		Crate:       crateHandler,
		Wastage:     wastageHandler,
		Expiry:      expiryHandler,
		Payment:     paymentHandler,
	}

	// Setup router with all handlers
	// Note: For E2E tests, we skip authentication to simplify testing
	// In production, routes are protected by middleware.Auth
	r := mux.NewRouter()

	// API v1 routes (no auth prefix for tests)
	api := r.PathPrefix("/api/v1").Subrouter()

	// Register all routes directly without auth middleware
	registerE2ERoutes(api, routerHandlers)

	// Create test server
	server := httptest.NewServer(r)

	return &E2ETestContext{
		DB:     db,
		Router: r,
		Server: server,
		Client: server.Client(),
	}
}

// registerE2ERoutes sets up all routes for E2E testing without authentication
func registerE2ERoutes(r *mux.Router, h *router.Handlers) {
	// Customer routes
	if h.Customer != nil {
		r.HandleFunc("/customers", h.Customer.Create).Methods("POST")
		r.HandleFunc("/customers", h.Customer.List).Methods("GET")
		r.HandleFunc("/customers/{id}", h.Customer.GetByID).Methods("GET")
		r.HandleFunc("/customers/{id}", h.Customer.Update).Methods("PUT")
		r.HandleFunc("/customers/{id}", h.Customer.Delete).Methods("DELETE")
	}

	// Inventory routes
	if h.Inventory != nil {
		r.HandleFunc("/inventory", h.Inventory.Create).Methods("POST")
		r.HandleFunc("/inventory", h.Inventory.List).Methods("GET")
		// Register specific routes BEFORE parameterized routes
		r.HandleFunc("/inventory/low-stock", h.Inventory.GetLowStock).Methods("GET")
		r.HandleFunc("/inventory/expiring", h.Inventory.GetExpiring).Methods("GET")
		// Parameterized routes come after
		r.HandleFunc("/inventory/{id}", h.Inventory.GetByID).Methods("GET")
		r.HandleFunc("/inventory/{id}", h.Inventory.Update).Methods("PUT")
		r.HandleFunc("/inventory/{id}", h.Inventory.Delete).Methods("DELETE")
		r.HandleFunc("/inventory/{id}/deduct", h.Inventory.DeductStock).Methods("POST")
	}

	// Transaction routes
	if h.Transaction != nil {
		r.HandleFunc("/transactions", h.Transaction.Create).Methods("POST")
		r.HandleFunc("/transactions", h.Transaction.List).Methods("GET")
		r.HandleFunc("/transactions/{id}", h.Transaction.GetByID).Methods("GET")
		r.HandleFunc("/transactions/{id}", h.Transaction.Update).Methods("PUT")
		r.HandleFunc("/transactions/{id}", h.Transaction.Delete).Methods("DELETE")
	}

	// SaleItem routes
	if h.SaleItem != nil {
		r.HandleFunc("/sale-items", h.SaleItem.Create).Methods("POST")
		r.HandleFunc("/sale-items/{id}", h.SaleItem.GetByID).Methods("GET")
		r.HandleFunc("/sale-items/{id}", h.SaleItem.Update).Methods("PUT")
		r.HandleFunc("/sale-items/{id}", h.SaleItem.Delete).Methods("DELETE")
		r.HandleFunc("/transactions/{transaction_id}/sale-items", h.SaleItem.ListByTransaction).Methods("GET")
	}

	// Crate routes
	if h.Crate != nil {
		r.HandleFunc("/crates/issue", h.Crate.IssueCrates).Methods("POST")
		r.HandleFunc("/crates/return", h.Crate.ReturnCrates).Methods("POST")
		r.HandleFunc("/crates/balance/{customerId}", h.Crate.GetBalance).Methods("GET")
		r.HandleFunc("/crates/history/{customerId}", h.Crate.GetHistory).Methods("GET")
		r.HandleFunc("/crates/{id}", h.Crate.GetByID).Methods("GET")
		r.HandleFunc("/crates/{id}", h.Crate.Update).Methods("PUT")
		r.HandleFunc("/crates/{id}", h.Crate.Delete).Methods("DELETE")
	}

	// Wastage routes
	if h.Wastage != nil {
		r.HandleFunc("/wastage", h.Wastage.RecordWastage).Methods("POST")
		r.HandleFunc("/wastage", h.Wastage.List).Methods("GET")
		r.HandleFunc("/wastage/{id}", h.Wastage.GetByID).Methods("GET")
		r.HandleFunc("/wastage/{id}", h.Wastage.Update).Methods("PUT")
		r.HandleFunc("/wastage/{id}", h.Wastage.Delete).Methods("DELETE")
		r.HandleFunc("/wastage/report", h.Wastage.GetReport).Methods("GET")
	}

	// Expiry alert routes
	if h.Expiry != nil {
		r.HandleFunc("/expiry-alerts/generate", h.Expiry.GenerateAlerts).Methods("POST")
		r.HandleFunc("/expiry-alerts", h.Expiry.List).Methods("GET")
		r.HandleFunc("/expiry-alerts/pending", h.Expiry.GetPending).Methods("GET")
		r.HandleFunc("/expiry-alerts/{id}", h.Expiry.GetByID).Methods("GET")
		r.HandleFunc("/expiry-alerts/{id}/acknowledge", h.Expiry.Acknowledge).Methods("POST")
		r.HandleFunc("/expiry-alerts/{id}", h.Expiry.Update).Methods("PUT")
		r.HandleFunc("/expiry-alerts/{id}", h.Expiry.Delete).Methods("DELETE")
	}

	// Payment routes
	if h.Payment != nil {
		r.HandleFunc("/payment-schedules", h.Payment.CreateSchedule).Methods("POST")
		r.HandleFunc("/payment-schedules/{id}", h.Payment.GetByID).Methods("GET")
		r.HandleFunc("/payment-schedules/customer/{customerId}", h.Payment.ListByCustomer).Methods("GET")
		r.HandleFunc("/payment-schedules/{id}/pay", h.Payment.RecordPayment).Methods("POST")
		r.HandleFunc("/payment-schedules/{id}", h.Payment.Update).Methods("PUT")
		r.HandleFunc("/payment-schedules/{id}", h.Payment.Delete).Methods("DELETE")
		r.HandleFunc("/payment-schedules/overdue/{customerId}", h.Payment.GetOverdue).Methods("GET")
	}
}

// TeardownE2ETest cleans up the test environment
func TeardownE2ETest(t *testing.T, ctx *E2ETestContext) {
	t.Helper()

	if ctx.Server != nil {
		ctx.Server.Close()
	}
	if ctx.DB != nil {
		TeardownTestDB(t, ctx.DB)
	}
}

// MakeRequest performs an HTTP request and returns the response
func (ctx *E2ETestContext) MakeRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		AssertNoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, ctx.Server.URL+path, bodyReader)
	AssertNoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ctx.Client.Do(req)
	AssertNoError(t, err)

	return resp
}

// DecodeResponse decodes the JSON response body
func (ctx *E2ETestContext) DecodeResponse(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
}

// ExtractData extracts the "data" field from API response wrapper
func (ctx *E2ETestContext) ExtractData(t *testing.T, resp *http.Response) map[string]interface{} {
	t.Helper()
	var wrapper map[string]interface{}
	ctx.DecodeResponse(t, resp, &wrapper)

	data, ok := wrapper["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Response does not contain data field. Got: %v", wrapper)
	}
	return data
}

// AssertHTTPStatusCode checks the HTTP status code for http.Response
func AssertHTTPStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()

	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("Expected status code %d, got %d. Body: %s", expected, resp.StatusCode, string(body))
	}
}

// ReadResponseBody reads and returns the response body as string
func ReadResponseBody(t *testing.T, resp *http.Response) string {
	t.Helper()

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	AssertNoError(t, err)
	return string(body)
}
