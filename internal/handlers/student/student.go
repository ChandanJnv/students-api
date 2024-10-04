package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/ChandanJnv/students-api/internal/types"
	"github.com/ChandanJnv/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
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
			slog.Info("Validation failed: ", err, student)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(err.(validator.ValidationErrors)))
			return
		}

		slog.Info("Created a student")

		response.WriteJSON(w, http.StatusCreated, student)
	}
}
