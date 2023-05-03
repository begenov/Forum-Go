package sqlite

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const path = "./schema/init.up.sql"

func InitDatabase(driver string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := initTable(path, db); err != nil {
		return nil, err
	}

	return db, nil
}

func initTable(path string, db *sql.DB) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, err = db.Exec(string(data))
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
