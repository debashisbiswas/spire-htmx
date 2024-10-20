package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
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

var templates = template.Must(template.ParseFiles("templates/index.html"))

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	err := templates.ExecuteTemplate(w, "index.html", entries)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
