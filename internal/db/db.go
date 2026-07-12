package db

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dbURL string) *sql.DB {
	database, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}

	if err := database.Ping(); err != nil {
		log.Fatal("Failed to Ping DB:", err)
	}
	log.Println("Connected to Neon Postgres successfully")

	createTable(database)

	return database
}

func createTable(databade *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		grade INT NOT NULL
	)`

	_, err := databade.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table: ", err)
	}
	log.Println("Student table ready")
}