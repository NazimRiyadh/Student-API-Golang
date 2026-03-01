package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/NazimRiyadh/student_api_golang/internal/storage"
	"github.com/NazimRiyadh/student_api_golang/internal/types"
	"github.com/NazimRiyadh/student_api_golang/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")

		var Student types.Student

		err := json.NewDecoder(r.Body).Decode(&Student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation

		if err := validator.New().Struct(Student); err != nil {
			validateErr := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValdationError(validateErr))
			return
		}

		lastid, err := storage.CreateStudent(Student.Name, Student.Email, Student.Age)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
		}
		slog.Info("Created User", slog.String("id", string(lastid)))

		response.WriteJSON(w, http.StatusCreated, map[string]int64{
			"id": lastid,
		})
	}
}
