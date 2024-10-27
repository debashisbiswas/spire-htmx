package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"spire/entry"
	"strconv"
	"strings"
	"time"

	libsqlvector "github.com/ryanskidmore/libsql-vector-go"
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
		)
	`)
	if err != nil {
		return err
	}

	// Everything breaks when I use the index. I suspect there's a bug in the
	// libsql-go client, or in the library I'm using to serialize embeddings
	// for the database. Leaving it out for now.

	_, err = db.Exec("CREATE INDEX entries_idx ON entries (libsql_vector_idx(embedding))")
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

	queryTemplate := fmt.Sprintf(
		"INSERT INTO entries (time, content, embedding) VALUES (?, ?, %s);",
		serializeEmbeddingsWithVectorPrefix(entry.Embedding.Slice()),
	)

	_, err = db.Exec(queryTemplate, entry.Time, entry.Content)
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

	rows, err := db.Query("SELECT time, content, vector_extract(embedding) FROM entries ORDER BY time DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entry.Entry

	for rows.Next() {
		var entry entry.Entry

		var timeString string
		var embeddingString string
		err := rows.Scan(&timeString, &entry.Content, &embeddingString)
		if err != nil {
			return nil, err
		}

		entry.Time, err = time.Parse("2006-01-02T15:04:05.999999999-07:00", timeString)
		if err != nil {
			log.Printf("Error parsing timestamp: %v\n", err)
			return nil, err
		}

		embedding, err := deserializeEmbeddings(embeddingString)
		if err != nil {
			log.Printf("Error parsing embeddings: %v\n", err)
			return nil, err
		}

		entry.Embedding = libsqlvector.NewVector(embedding)

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
		SELECT time, content, vector_extract(embedding)
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
		var embeddingString string
		err := rows.Scan(&timeString, &entry.Content, &embeddingString)
		if err != nil {
			return nil, err
		}

		entry.Time, err = time.Parse("2006-01-02T15:04:05.999999999-07:00", timeString)
		if err != nil {
			log.Printf("Error parsing timestamp: %v\n", err)
			return nil, err
		}

		embedding, err := deserializeEmbeddings(embeddingString)
		if err != nil {
			log.Printf("Error parsing embeddings: %v\n", err)
			return nil, err
		}

		entry.Embedding = libsqlvector.NewVector(embedding)

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *SQLiteStorage) SearchEntriesEmbedding(embedding libsqlvector.Vector) ([]entry.Entry, error) {
	db, err := s.getDatabaseConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := fmt.Sprintf(`
		SELECT time, content, vector_extract(embedding)
		FROM entries
		ORDER BY vector_distance_cos(embedding, vector('[%s]'))
	`, serializeEmbeddings(embedding.Slice()))

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []entry.Entry

	for rows.Next() {
		var entry entry.Entry

		var timeString string
		var embeddingStr string
		err := rows.Scan(&timeString, &entry.Content, &embeddingStr)
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

// floats [1, 2, 3] -> string '[1,2,3]'
func serializeEmbeddings(input []float32) string {
	strNums := make([]string, len(input))
	for i, num := range input {
		strNums[i] = strconv.FormatFloat(float64(num), 'f', -1, 32)
	}

	builder := strings.Builder{}
	builder.WriteString("'[")
	builder.WriteString(strings.Join(strNums, ","))
	builder.WriteString("]'")

	return builder.String()
}

// string [1,2,3] -> floats [1, 2, 3]
func deserializeEmbeddings(input string) ([]float32, error) {
	var result []float32

	log.Println(input)

	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// floats [1, 2, 3] -> string vector('[1,2,3]')
func serializeEmbeddingsWithVectorPrefix(input []float32) string {
	builder := strings.Builder{}
	builder.WriteString("vector(")
	builder.WriteString(serializeEmbeddings(input))
	builder.WriteString(")")

	return builder.String()
}
