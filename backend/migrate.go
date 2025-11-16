
package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4"
)

func runMigrations(migrationsPath, databaseURL string) error {
    if migrationsPath == "" {
        migrationsPath = "file://infra/migrations"
    }
    if databaseURL == "" {
        return nil
    }
    m, err := migrate.New(migrationsPath, databaseURL)
    if err != nil {
        return err
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return err
    }
    return nil
}
