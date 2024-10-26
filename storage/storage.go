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
			content TEXT NOT NULL,
			embedding F32_BLOB(512)
		);

		CREATE INDEX entries_idx ON entries (libsql_vector_idx(embedding));
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

	_, err = db.Exec(`
		INSERT INTO entries (time, content, embedding) VALUES (?, ?, ?);
	`, entry.Time, entry.Content, entry.Embedding)
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

	rows, err := db.Query("SELECT time, content, embedding FROM entries ORDER BY time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entry.Entry

	for rows.Next() {
		var entry entry.Entry

		var timeString string
		err := rows.Scan(&timeString, &entry.Content, &entry.Embedding)
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

func (s *SQLiteStorage) SearchEntries(query string) ([]entry.Entry, error) {
	// TODO: This is repeated quite a bit. Is there a better way, maybe
	// something similar to Python's ContextManager?
	db, err := s.getDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// TODO: I don't really need the embedding here, but I'm getting it because
	// my test uses it, and it doesn't feel right to leave the struct field
	// empty.
	rows, err := db.Query(`
		SELECT time, content, embedding
		FROM entries
		WHERE content LIKE ?
		ORDER BY time DESC
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entry.Entry

	for rows.Next() {
		var entry entry.Entry

		var timeString string
		err := rows.Scan(&timeString, &entry.Content, &entry.Embedding)
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
