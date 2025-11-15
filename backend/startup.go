
package main

import (
    "log"
    "os"
    "time"
    "context"
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func initApp() {
    // Load .env if present
    _ = os.Setenv("MIGRATE_ON_START", os.Getenv("MIGRATE_ON_START"))
}

func prepareDB() (*sql.DB, error) {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        return nil, nil
    }
    // optionally run migrations
    if os.Getenv("MIGRATE_ON_START") == "true" {
        log.Println("Running DB migrations...")
        if err := runMigrations("file://infra/migrations", dsn); err != nil {
            log.Printf("migration error: %v", err)
            // continue - do not fail startup hard; you may change behaviour
        }
    }
    db, err := sql.Open("pgx", dsn)
    if err != nil { return nil, err }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := db.PingContext(ctx); err != nil { return nil, err }
    return db, nil
}
