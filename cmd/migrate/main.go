package main

import (
	"fmt"
	"log"
	"os"

	"src/internal/database"
	_ "src/migrations"

	"github.com/pressly/goose/v3"
)

func main() {
	// Initialize DB (loads env via godotenv autoload in database package)
	_ = database.New()
	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("goose dialect error: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatalf("usage: migrate [up|down|status]")
	}

	db := database.SQLDB()
	dir := "migrations"
	var err error

	switch os.Args[1] {
	case "up":
		err = goose.Up(db, dir)
	case "down":
		err = goose.Down(db, dir)
	case "status":
		err = goose.Status(db, dir)
	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		log.Fatalf("migration command failed: %v", err)
	}
	fmt.Println("migration command ok")
}
