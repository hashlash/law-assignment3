package main

import (
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
	"os"
)

func home(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./pages/home.html")
}

func upload(w http.ResponseWriter, r *http.Request) {
	routingKey := uuid.New()

	client := &http.Client{}
	url := "http://" + os.Getenv("SERVER2_HOST")

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		log.Println("Failed to create request object:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("X-ROUTING-KEY", routingKey.String())
	req.Header.Add("Content-Type", r.Header.Get("Content-Type"))
	req.Header.Add("Content-Length", r.Header.Get("Content-Length"))

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed to do the request:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	fmt.Println(resp)

	tmpl, err := template.ParseFiles("./pages/progress.html")
	if err != nil {
		log.Println("Failed to parse html file:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, map[string]string {"routingKey": routingKey.String()})
	if err != nil {
		log.Println("Failed to execute template:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}
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
	log.Println("Serving on:", os.Getenv("SERVER1_HOST"))
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER1_HOST"), nil))
}
