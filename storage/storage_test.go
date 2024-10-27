package storage

import (
	"os"
	"reflect"
	"spire/entry"
	"testing"
	"time"

	libsqlvector "github.com/ryanskidmore/libsql-vector-go"
	"golang.org/x/exp/rand"
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

	randomEmbeddings := generateRandomEmbeddings()

	originalEntries := []entry.Entry{
		{
			Time:      time.Now(),
			Content:   "welcome to the playground",
			Embedding: libsqlvector.NewVector(randomEmbeddings),
		},
		{
			Time:      time.Now(),
			Content:   "follow me",
			Embedding: libsqlvector.NewVector(randomEmbeddings),
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

	foundEmbedding := searchResult[0].Embedding.Slice()
	expectedEmbedding := randomEmbeddings

	if len(foundEmbedding) != len(expectedEmbedding) {
		t.Errorf("embedding length is incorrect: found %d, expected %d", len(foundEmbedding), len(expectedEmbedding))
		t.FailNow()
	}

	if !reflect.DeepEqual(foundEmbedding, expectedEmbedding) {
		t.Errorf("embeddings are %v\nexpected %v", foundEmbedding, expectedEmbedding)
		t.FailNow()
	}
}

func generateRandomEmbeddings() []float32 {
	rand.Seed(42)

	result := make([]float32, 512)

	for i := range result {
		result[i] = rand.Float32()
	}

	return result
}
