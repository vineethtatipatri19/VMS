package router

import (
	"net/http"

	"github.com/example/pgvms/internal/handlers"
	"github.com/example/pgvms/internal/middleware"
	"github.com/gorilla/mux"
)

// Config holds the configuration for the router
type Config struct {
	JWTSecret []byte
	Handlers  *Handlers
}

// Handlers holds all HTTP handlers
type Handlers struct {
	Customer    *handlers.CustomerHandler
	Inventory   *handlers.InventoryHandler
	Transaction *handlers.TransactionHandler
	SaleItem    *handlers.SaleItemHandler
	Crate       *handlers.CrateHandler
	Wastage     *handlers.WastageHandler
	Expiry      *handlers.ExpiryHandler
	Payment     *handlers.PaymentHandler

	// Auth handlers (from old main.go)
	Register http.HandlerFunc
	Login    http.HandlerFunc
	Health   http.HandlerFunc

	// Legacy handlers (to be migrated)
	Dashboard       http.HandlerFunc
	RecentActivity  http.HandlerFunc
	OverduePayments http.HandlerFunc
	WastageSummary  http.HandlerFunc
	Forecast        http.HandlerFunc
	GenerateReport  http.HandlerFunc
}

// Setup creates and configures the main router with all routes and middleware
func Setup(cfg Config) *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.Logging)
	r.Use(middleware.CORS)

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Public routes
	registerPublicRoutes(api, cfg.Handlers)

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))
	registerProtectedRoutes(protected, cfg.Handlers)

	return r
}

// registerPublicRoutes sets up routes that don't require authentication
func registerPublicRoutes(r *mux.Router, h *Handlers) {
	r.HandleFunc("/register", h.Register).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/login", h.Login).Methods(http.MethodPost, http.MethodOptions)
} // registerProtectedRoutes sets up routes that require authentication
func registerProtectedRoutes(r *mux.Router, h *Handlers) {
	// Health check
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)

	// Dashboard
	r.HandleFunc("/dashboard", h.Dashboard).Methods(http.MethodGet)
	r.HandleFunc("/dashboard/activity", h.RecentActivity).Methods(http.MethodGet)

	// Register resource-specific routes
	registerCustomerRoutes(r, h.Customer)
	registerInventoryRoutes(r, h.Inventory)
	registerTransactionRoutes(r, h.Transaction)
	registerSaleItemRoutes(r, h.SaleItem)
	registerCrateRoutes(r, h.Crate)
	registerWastageRoutes(r, h.Wastage)
	registerExpiryRoutes(r, h.Expiry)
	registerPaymentRoutes(r, h.Payment)

	// Reports
	r.HandleFunc("/reports/overdue-payments", h.OverduePayments).Methods(http.MethodGet)
	r.HandleFunc("/reports/wastage-summary", h.WastageSummary).Methods(http.MethodGet)
	r.HandleFunc("/reports/generate", h.GenerateReport).Methods(http.MethodPost)

	// Forecasting (AI)
	r.HandleFunc("/forecast", h.Forecast).Methods(http.MethodPost)
}
