package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Entry struct {
	Time    time.Time
	Content string
}

var entries = []Entry{
	{time.Now(), "welcome to the playground"},
	{time.Now(), "follow me"},
}

var templates = template.Must(template.ParseFiles(
	"templates/index.html",
	"templates/components/entry.html",
))

func baseHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	err := templates.ExecuteTemplate(w, "index.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func newEntryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	content := r.PostForm.Get("entry")
	newEntry := Entry{time.Now(), content}
	entries = append(entries, newEntry)

	err := templates.ExecuteTemplate(w, "entry.html", newEntry)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("GET /", baseHandler)
	http.HandleFunc("POST /entries", newEntryHandler)

	port := 8080
	portString := strconv.Itoa(port)
	log.Println("Listening on port " + portString)
	log.Fatal(http.ListenAndServe(":"+portString, nil))
}
