package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/NazimRiyadh/student_api_golang/internal/config"
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

	res, err := db.Exec(`CREATE TABLE IF NOT EXISTS STUDENTS (
		ID INTEGER PRIMARY KEY AUTOINCREMENT,
		NAME TEXT,
		EMAIL TEXT,
		AGE INTEGER
	)`)
	fmt.Println(res) //remove it later
	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (r *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {

	stmt, err := r.Db.Prepare("INSERT INTO STUDENTS (NAME, EMAIL, AGE) VALUES (?, ?, ?)")
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
