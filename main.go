package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"spire/entry"
	"spire/storage"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/components/entry.html",
	"templates/components/entries.html",
))

type Server struct {
	Storage storage.SQLiteStorage
}

func (server *Server) baseHandler(w http.ResponseWriter, r *http.Request) {
	entries, err := server.Storage.GetEntries()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "index.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) newEntryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	content := r.PostForm.Get("entry")
	newEntry := entry.Entry{Time: time.Now(), Content: content}
	err := server.Storage.SaveEntry(newEntry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "entry.html", newEntry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) searchHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	content := r.PostForm.Get("search")

	var entries []entry.Entry
	var err error

	if content == "" {
		entries, err = server.Storage.GetEntries()
	} else {
		entries, err = server.Storage.SearchEntries(content)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = templates.ExecuteTemplate(w, "entries.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env")
	}

	store, err := storage.NewSQLiteStorage("main.db")
	if err != nil {
		log.Fatalf("error initializing database: %v\n", err)
	}

	server := Server{*store}

	http.HandleFunc("GET /", server.baseHandler)
	http.HandleFunc("POST /entries", server.newEntryHandler)
	http.HandleFunc("POST /search", server.searchHandler)

	port := 8080
	portString := strconv.Itoa(port)
	log.Println("Listening on port " + portString)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
