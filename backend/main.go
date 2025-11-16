
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

func main() {
	// Load JWT secret from environment
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		jwtKey = []byte(jwtSecret)
	}
	
	dsn := os.Getenv("DATABASE_URL") // expected: postgres://user:pass@host:5432/dbname
	if dsn == "" {
		log.Fatal("DATABASE_URL env var required")
	}
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("db open: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	// Run migrations if enabled
	if os.Getenv("MIGRATE_ON_START") == "true" {
		log.Println("Running DB migrations...")
		if err := runMigrations("file://infra/migrations", dsn); err != nil {
			log.Printf("migration warning: %v", err)
		}
	}

	r := chi.NewRouter()
	
	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)
	
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		
		// Protected routes
		r.Group(func(rr chi.Router) {
			rr.Use(authMiddleware)
			
			// Health check
			rr.Get("/health", healthHandler)
			
			// Dashboard
			rr.Get("/dashboard", dashboardHandler)
			rr.Get("/dashboard/activity", recentActivityHandler)
			
			// Inventory
			rr.Get("/inventory", listInventory)
			rr.Post("/inventory", createInventory)
			rr.Get("/inventory/{id}", getInventoryItem)
			rr.Put("/inventory/{id}", updateInventory)
			rr.Delete("/inventory/{id}", deleteInventory)
			
			// Customers
			rr.Get("/customers", listCustomers)
			rr.Post("/customers", createCustomer)
			rr.Get("/customers/{id}", getCustomer)
			rr.Put("/customers/{id}", updateCustomer)
			rr.Delete("/customers/{id}", deleteCustomer)
			
			// Transactions
			rr.Get("/transactions", listTransactions)
			rr.Post("/transactions", createTransaction)
			rr.Get("/transactions/{id}", getTransaction)
			
			// Crates
			rr.Get("/crates", listCrates)
			rr.Post("/crates", createCrateEntry)
			rr.Get("/crates/balance/{customerId}", getCrateBalance)
			
			// Forecasting (AI)
			rr.Post("/forecast", forecastHandler)
			
			// Reports
			rr.Post("/reports/generate", generateReportHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	http.ListenAndServe(":"+port, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// CORS middleware to allow frontend access
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

