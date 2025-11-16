package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/example/pgvms/internal/config"
	"github.com/example/pgvms/internal/handlers"
	"github.com/example/pgvms/internal/repository/postgres"
	"github.com/example/pgvms/internal/router"
	"github.com/example/pgvms/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Global DB for legacy handlers
var db *sql.DB

func main() {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	db, err = sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Ping database to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Database connection established")

	// Run migrations if enabled
	if cfg.Database.MigrateOnStart {
		log.Println("Running DB migrations...")
		if err := runMigrations(cfg.Database.MigrationsPath, cfg.Database.URL); err != nil {
			log.Printf("Migration warning: %v", err)
		}
	}

	// Initialize repositories
	customerRepo := postgres.NewCustomerRepository(db)
	inventoryRepo := postgres.NewInventoryRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	saleItemRepo := postgres.NewSaleItemRepository(db)
	crateRepo := postgres.NewCrateRepository(db)
	wastageRepo := postgres.NewWastageRepository(db)
	expiryRepo := postgres.NewExpiryAlertRepository(db)
	paymentRepo := postgres.NewPaymentScheduleRepository(db)

	// Initialize services with repository dependencies
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

	// Initialize handlers with service dependencies
	customerHandler := handlers.NewCustomerHandler(customerService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	saleItemHandler := handlers.NewSaleItemHandler(saleItemService)
	crateHandler := handlers.NewCrateHandler(crateService)
	wastageHandler := handlers.NewWastageHandler(wastageService)
	expiryHandler := handlers.NewExpiryHandler(expiryService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	// Setup router with all handlers
	routerConfig := router.Config{
		JWTSecret: cfg.JWT.Secret,
		Handlers: &router.Handlers{
			Customer:    customerHandler,
			Inventory:   inventoryHandler,
			Transaction: transactionHandler,
			SaleItem:    saleItemHandler,
			Crate:       crateHandler,
			Wastage:     wastageHandler,
			Expiry:      expiryHandler,
			Payment:     paymentHandler,

			// Auth handlers
			Register: registerHandler,
			Login:    loginHandler,
			Health:   healthHandler,

			// Legacy handlers (to be refactored)
			Dashboard:       dashboardHandler,
			RecentActivity:  recentActivityHandler,
			OverduePayments: getOverduePayments,
			WastageSummary:  getWastageSummary,
			Forecast:        forecastHandler,
			GenerateReport:  generateReportHandler,
		},
	}

	r := router.Setup(routerConfig)

	// Start server
	port := cfg.Server.Port
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
