package storage

import "github.com/ChandanJnv/students-api/internal/types"

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
	GetById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	DeleteById(id int64) error
}
