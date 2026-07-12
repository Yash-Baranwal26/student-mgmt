package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Yash-Baranwal26/student-mgmt/internal/models"
	"github.com/Yash-Baranwal26/student-mgmt/internal/repository"
)

type StudentHandler struct {
	Repo *repository.StudentRepository
}

func NewStudentHandler(repo *repository.StudentRepository) *StudentHandler {
	return &StudentHandler{Repo: repo}
}

func (h *StudentHandler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var s models.Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.Repo.Create(&s); err != nil {
		http.Error(w, "Failed to insert student: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func (h *StudentHandler) GetStudents(w http.ResponseWriter, r *http.Request) {
	students, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "Failed to fetch students: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}