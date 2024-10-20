package main

import (
	"html/template"
	"log"
	"net/http"
	"spire/entry"
	"spire/storage"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/components/entry.html",
	"templates/components/entries.html",
))

func baseHandler(w http.ResponseWriter, r *http.Request, store *storage.SQLiteStorage) {
	entries, err := store.GetEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "index.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newEntryHandler(w http.ResponseWriter, r *http.Request, store *storage.SQLiteStorage) {
	r.ParseForm()

	content := r.PostForm.Get("entry")
	newEntry := entry.Entry{Time: time.Now(), Content: content}
	err := store.SaveEntry(newEntry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "entry.html", newEntry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func searchHandler(w http.ResponseWriter, r *http.Request, store *storage.SQLiteStorage) {
	r.ParseForm()

	content := r.PostForm.Get("search")

	var entries []entry.Entry
	var err error

	if content == "" {
		entries, err = store.GetEntries()
	} else {
		entries, err = store.SearchEntries(content)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "entries.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handlerWithStorage(
	fn func(http.ResponseWriter, *http.Request, *storage.SQLiteStorage),
	store *storage.SQLiteStorage,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, store)
	}
}

func main() {
	store, err := storage.NewSQLiteStorage("main.db")
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
	}

	http.HandleFunc("GET /", handlerWithStorage(baseHandler, store))
	http.HandleFunc("POST /entries", handlerWithStorage(newEntryHandler, store))
	http.HandleFunc("POST /search", handlerWithStorage(searchHandler, store))

	port := 8080
	portString := strconv.Itoa(port)
	log.Println("Listening on port " + portString)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
