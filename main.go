package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// Student represents a row in our students table
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Grade int    `json:"grade"`
}

var db *sql.DB

func main() {
	// Load .env file into environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system env vars")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Open connection (pgx registers itself as "pgx" driver via stdlib import above)
	var err error
	db, err = sql.Open("pgx", dbURL)
	if err != nil {
		log.Fatal("Failed to open DB connection:", err)
	}
	defer db.Close()

	// Actually test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}
	log.Println("Connected to Neon Postgres successfully")

	createTable()

	r := chi.NewRouter()
	r.Post("/students", createStudent)
	r.Get("/students", getStudents)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// createTable creates the students table if it doesn't already exist
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		grade INT NOT NULL
	);`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("students table ready")
}

// createStudent handles POST /students
func createStudent(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO students (name, email, grade) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, s.Name, s.Email, s.Grade).Scan(&s.ID)
	if err != nil {
		http.Error(w, "Failed to insert student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// getStudents handles GET /students
func getStudents(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, name, email, grade FROM students`)
	if err != nil {
		http.Error(w, "Failed to fetch students: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Grade); err != nil {
			http.Error(w, "Failed to scan student: "+err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}