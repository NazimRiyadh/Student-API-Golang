package student

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/NazimRiyadh/student_api_golang/internal/storage"
	"github.com/NazimRiyadh/student_api_golang/internal/types"
	"github.com/NazimRiyadh/student_api_golang/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		slog.Info("Creating a student")

		var Student types.Student

		defer r.Body.Close()

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
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		slog.Info("Created User", slog.Int64("id", lastid))

		response.WriteJSON(w, http.StatusCreated, map[string]int64{
			"id": lastid,
		})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("getting student", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			slog.Error("Convertion of id to int failed ", slog.String("id", id))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("Error getting user", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")
		students, err := storage.GetStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, err)
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}
