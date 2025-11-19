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
	Auth        *handlers.AuthHandler
	Dashboard   *handlers.DashboardHandler
	Report      *handlers.ReportHandler
	Forecast    *handlers.ForecastHandler

	// Simple handlers
	Health http.HandlerFunc
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
	r.HandleFunc("/register", h.Auth.Register).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/login", h.Auth.Login).Methods(http.MethodPost, http.MethodOptions)
} // registerProtectedRoutes sets up routes that require authentication
func registerProtectedRoutes(r *mux.Router, h *Handlers) {
	// Health check
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)

	// Dashboard
	r.HandleFunc("/dashboard", h.Dashboard.GetStats).Methods(http.MethodGet)
	r.HandleFunc("/dashboard/activity", h.Dashboard.GetRecentActivity).Methods(http.MethodGet)

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
	r.HandleFunc("/reports/generate", h.Report.GenerateReport).Methods(http.MethodPost)

	// Forecasting (AI)
	r.HandleFunc("/forecast", h.Forecast.GenerateForecast).Methods(http.MethodPost)
}
