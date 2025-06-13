package sqlite

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

// DB INIT
func Connect() (*Database, error) {
	var err error
	db, err := sql.Open("sqlite3", "database/database.db")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	err = applyMigrations(db)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
