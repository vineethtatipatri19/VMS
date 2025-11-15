
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
)

var db *sql.DB

func main() {
	dsn := os.Getenv("DATABASE_URL") // expected: postgres://user:pass@host:5432/dbname
	if dsn == "" {
		log.Fatal("DATABASE_URL env var required")
	}
	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil { log.Fatalf("db open: %v", err) }
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil { log.Fatalf("db ping: %v", err) }

	r := chi.NewRouter()
	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/register", registerHandler)
		r.Post("/login", loginHandler)
		// protected routes
		r.Group(func(rr chi.Router){
			rr.Use(authMiddleware)
			rr.Get("/health", healthHandler)
			rr.Get("/inventory", listInventory)
			rr.Post("/inventory", createInventory)
			rr.Get("/customers", listCustomers)
			rr.Post("/customers", createCustomer)
			rr.Get("/transactions", listTransactions)
			rr.Post("/transactions", createTransaction)
			rr.Post("/forecast", forecastHandler) // AI stub
		})
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	log.Printf("listening on :%s", port)
	http.ListenAndServe(":"+port, r)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
}
