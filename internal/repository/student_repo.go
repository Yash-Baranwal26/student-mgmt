package repository

import (
	"database/sql"

	"github.com/Yash-Baranwal26/student-mgmt/internal/models"
)

type StudentRepository struct {
	DB *sql.DB
}

func NewStudentRepository(db *sql.DB) *StudentRepository {
	return &StudentRepository{DB: db}
}

func (r *StudentRepository) Create(s *models.Student) error {
	query := `INSERT INTO Students (name, email, grade) Values ($1, $2, $3) RETURNING id`
	return r.DB.QueryRow(query, s.Name, s.Email, s.Grade).Scan(&s.ID)
}

func (r *StudentRepository) GetAll() ([]models.Student, error) {
	rows, err:= r.DB.Query(`SELECT id, name, email, grade FROM students`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Grade); err != nil {
			return nil, err
		}
		students = append(students, s)
	}
	return students, nil
}