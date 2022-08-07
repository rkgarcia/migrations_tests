package main

import (
	"embed"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	"log"
)

//go:embed migrations/*.sql
var migrations embed.FS

func main() {
	db, err := sqlx.Connect("sqlserver", "sqlserver://localhost?database=migrations_test")
	if err != nil {
		log.Fatal(err)
	}
	driver, err := sqlserver.WithInstance(db.DB, &sqlserver.Config{})
	if err != nil {
		log.Fatal(err)
	}
	d, err := iofs.New(migrations, "migrations")
	m, err := migrate.NewWithInstance("iofs", d, "sqlserver", driver)
	if err != nil {
		log.Fatal(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("We can't migrate database error: %v", err)
	}
	v, dirty, err := m.Version()
	if err != nil {
		log.Printf("Error getting version: %v", err)
	}
	log.Printf("Current database version: %d dirty? %v", v, dirty)
}
