package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ChandanJnv/students-api/internal/config"
	"github.com/ChandanJnv/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	if _, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		age INTEGER NOT NULL
	)`); err != nil {
		return nil, err
	}

	return &sqlite{Db: db}, nil
}

func (s *sqlite) CreateStudent(name, email string, age int) (int64, error) {

	stmt, err := s.Db.Prepare(`INSERT INTO students (name, email, age) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return lastId, err
	}

	return lastId, nil
}

func (s *sqlite) GetById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare(`SELECT id, name, email, age FROM students WHERE id = ? LIMIT 1`)
	if err != nil {
		slog.Error("failed to prepare database statement", slog.String("error", err.Error()))
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	if err := stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
		slog.Error("failed to query database", slog.String("error", err.Error()))
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("student with id %d not found: %w", id, err)
		}
		return types.Student{}, fmt.Errorf("querry error: %w", err)
	}

	return student, nil
}

func (s *sqlite) GetAllStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare(`SELECT id, name, email, age FROM students`)
	if err != nil {
		slog.Error("failed to prepare database statement", slog.String("error", err.Error()))
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		slog.Error("failed to query database", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	var students []types.Student
	for rows.Next() {
		var student types.Student
		if err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age); err != nil {
			slog.Error("failed to scan database row", slog.String("error", err.Error()))
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil
}

func (s *sqlite) DeleteById(id int64) error {
	_, err := s.GetById(id)
	if err != nil {
		slog.Error("student not found", slog.String("error", err.Error()))
		return err
	}

	stmt, err := s.Db.Prepare(`DELETE FROM students WHERE id = ?`)
	if err != nil {
		slog.Error("failed to prepare database statement", slog.String("error", err.Error()))
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(id); err != nil {
		slog.Error("failed to delete student", slog.String("error", err.Error()))
		return err
	}
	slog.Info("deleted student", slog.Int64("id", id))

	return nil
}
