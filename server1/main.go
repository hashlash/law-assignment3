package main

import (
	"fmt"
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/home.html")
}

func upload(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		home(w, r)
	} else if r.Method == http.MethodPost {
		upload(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func setupRoutes() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8888", nil))
}

func main() {
	setupRoutes()
}
