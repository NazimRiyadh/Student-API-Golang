package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/NazimRiyadh/student_api_golang/internal/config"
	"github.com/NazimRiyadh/student_api_golang/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME TEXT,
		EMAIL TEXT,
		AGE INTEGER
	)`)
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (r *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := r.Db.Prepare("INSERT INTO STUDENTS (NAME, EMAIL, AGE) VALUES ( ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastid, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (r *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := r.Db.Prepare("SELECT ID, NAME, EMAIL, AGE FROM STUDENTS WHERE ID=? LIMIT 1")
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("No student found with id: %s", fmt.Sprint(id))
		}

		return types.Student{}, fmt.Errorf("query error: %s", err)
	}
	return student, nil
}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT ID, NAME, EMAIL, AGE FROM STUDENTS")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}
