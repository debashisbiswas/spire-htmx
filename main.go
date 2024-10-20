package main

import (
	"fmt"
	"log"
	"net/http"
)

func baseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	fmt.Fprintln(w, "Welcome to Spire!")
}

func main() {
	http.HandleFunc("/", baseHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
