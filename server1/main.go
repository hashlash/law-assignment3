package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/home.html")
}

func upload(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	url := "http://" + os.Getenv("SERVER2_HOST")

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		log.Println("Failed to create request object", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to do the request", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(resp)

	http.ServeFile(w, r, "./pages/progress.html")
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
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER1_HOST"), nil))
}
