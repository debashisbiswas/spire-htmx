package storage

import (
	"os"
	"spire/entry"
	"testing"
	"time"
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
		{Time: time.Now(), Content: "welcome to the playground"},
		{Time: time.Now(), Content: "follow me"},
	}

	for _, e := range originalEntries {
		store.SaveEntry(e)
	}

	entries, err := store.GetEntries()
	if err != nil {
		t.Errorf("error getting entries: %v\n", err)
	}

	if len(entries) != len(originalEntries) {
		t.Errorf("found %d entries, but original had %d", len(entries), len(originalEntries))
	}

	searchResult, err := store.SearchEntries("playground")
	if err != nil {
		t.Errorf("error searching entries: %v\n", err)
	}

	if len(searchResult) != 1 {
		t.Errorf("search found %d entries, but expected %d", len(searchResult), 1)
	}
}
