package storage

import (
	"database/sql"
	"log"
	"spire/entry"
	"time"

	_ "github.com/tursodatabase/go-libsql"
)

type SQLiteStorage struct {
	databaseName string
}

func NewSQLiteStorage(name string) (*SQLiteStorage, error) {
	storage := &SQLiteStorage{databaseName: name}
	err := storage.init()

	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *SQLiteStorage) getDatabaseConnection() (*sql.DB, error) {
	db, err := sql.Open("libsql", "file:"+s.databaseName)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func (s *SQLiteStorage) init() error {
	db, err := s.getDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS entries (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time TIMESTAMP NOT NULL,
			content TEXT NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) SaveEntry(entry entry.Entry) error {
	db, err := s.getDatabaseConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO entries (time, content) VALUES (?, ?)", entry.Time, entry.Content)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) GetEntries() ([]entry.Entry, error) {
	db, err := s.getDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT time, content FROM entries ORDER BY time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entry.Entry

	for rows.Next() {
		var entry entry.Entry

		var timeString string
		err := rows.Scan(&timeString, &entry.Content)
		if err != nil {
			return nil, err
		}

		entry.Time, err = time.Parse("2006-01-02T15:04:05.999999999-07:00", timeString)
		if err != nil {
			log.Printf("Error parsing timestamp: %v\n", err)
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
