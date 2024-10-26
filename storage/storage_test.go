package storage

import (
	"os"
	"reflect"
	"spire/entry"
	"testing"
	"time"

	"github.com/ryanskidmore/libsql-vector-go"
)

func TestStorage(t *testing.T) {
	testDatabasePath := "storage_test.db"

	store, err := NewSQLiteStorage(testDatabasePath)

	if err != nil {
		t.Errorf("error creating storage: %v\n", err)
	}

	defer func() {
		os.Remove(testDatabasePath)
	}()

	originalEntries := []entry.Entry{
		{
			Time:      time.Now(),
			Content:   "welcome to the playground",
			Embedding: libsqlvector.NewVector([]float32{0.3, 0.6}),
		},
		{
			Time:      time.Now(),
			Content:   "follow me",
			Embedding: libsqlvector.NewVector([]float32{1, 2}),
		},
	}

	for _, e := range originalEntries {
		err = store.SaveEntry(e)
		if err != nil {
			t.Errorf("error saving entry: %v\n", err)
			t.FailNow()
		}
	}

	entries, err := store.GetEntries()
	if err != nil {
		t.Errorf("error getting entries: %v\n", err)
		t.FailNow()
	}

	if len(entries) != len(originalEntries) {
		t.Errorf("found %d entries, but original had %d", len(entries), len(originalEntries))
		t.FailNow()
	}

	searchResult, err := store.SearchEntries("playground")
	if err != nil {
		t.Errorf("error searching entries: %v\n", err)
		t.FailNow()
	}

	if len(searchResult) != 1 {
		t.Errorf("search found %d entries, but expected %d", len(searchResult), 1)
		t.FailNow()
	}

	foundEmbedding := searchResult[0].Embedding
	expectedEmbedding := libsqlvector.NewVector([]float32{0.3, 0.6})
	if !reflect.DeepEqual(foundEmbedding, expectedEmbedding) {
		t.Errorf("embeddings are %v, but expected %v", foundEmbedding, expectedEmbedding)
		t.FailNow()
	}
}
