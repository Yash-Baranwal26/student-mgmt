package router

import (
	"github.com/go-chi/chi/v5"

	"github.com/Yash-Baranwal26/student-mgmt/internal/handlers"
)

func New(studentHandler *handlers.StudentHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/students", studentHandler.CreateStudent)
	r.Get("/students", studentHandler.GetStudents)

	return r
}