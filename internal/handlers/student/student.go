package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/ChandanJnv/students-api/internal/storage"
	"github.com/ChandanJnv/students-api/internal/types"
	"github.com/ChandanJnv/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")
		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		} else if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validate request
		if err := validator.New().Struct(student); err != nil {
			slog.Info("Validation failed: ", slog.String("error", err.Error()), slog.String("name", student.Name))
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		studentID, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			slog.Info("Failed to create a student: ", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Created a student ", slog.Int64("id", studentID))

		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": studentID})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			slog.Error("failed to Get student by id", slog.String("id", id))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("id is required")))
			return
		}
		slog.Info("Get student by id", slog.String("id", id))

		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("failed to parse id into int64", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		studentDetail, err := storage.GetById((studentID))
		if err != nil {
			slog.Error("failed to Get student by id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, studentDetail)

	}
}

func GetAllStudents(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Get all students")
		students, err := storage.GetAllStudents()
		if err != nil {
			slog.Error("failed to Get all students", slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}

func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			slog.Error("failed to Get student by id", slog.String("id", id))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("id is required")))
			return
		}
		slog.Info("Get student by id", slog.String("id", id))

		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("failed to parse id into int64", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := storage.DeleteById((studentID)); err != nil {
			slog.Error("failed to Delete student by id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, map[string]string{"result": "successs"})

	}
}

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			slog.Error("failed to Get student by id", slog.String("id", id))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(errors.New("id is required")))
			return
		}

		studentID, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("failed to parse id into int64", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		slog.Info("updating student by id", slog.String("id", id))
		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		} else if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// validate request
		if err := validator.New().Struct(student); err != nil {
			slog.Info("Validation failed: ", slog.String("error", err.Error()), slog.String("name", student.Name))
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		if err := storage.UpdateById(studentID, student.Name, student.Email, student.Age); err != nil {
			slog.Error("failed to update student by id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, map[string]string{"result": "successs"})
	}
}
