package main

import (
	"log"
	"net/http"

	"github.com/Yash-Baranwal26/student-mgmt/internal/config"
	"github.com/Yash-Baranwal26/student-mgmt/internal/db"
	"github.com/Yash-Baranwal26/student-mgmt/internal/handlers"
	"github.com/Yash-Baranwal26/student-mgmt/internal/repository"
	"github.com/Yash-Baranwal26/student-mgmt/internal/router"
)


func main() {
	cfg := config.Load()

	database := db.Connect(cfg.DatabaseURL)
	defer database.Close()

	studentRepo := repository.NewStudentRepository(database)
	studentHandler := handlers.NewStudentHandler(studentRepo)

	r := router.New(studentHandler)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}