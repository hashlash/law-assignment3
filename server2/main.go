package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"

	//"net/http"
	"os"
)

//func failOnError(w http.ResponseWriter, err error, msg string) {
func failOnError(err error, msg string) {
	if err != nil {
		//http.Error(w, msg, 500)
		log.Fatalf("%s: %s", msg, err)
	}
}

//func rmqSend(w http.ResponseWriter, msg string) {
func rmqSend() {
	time.Sleep(3 * time.Second)
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Println("Failed to connect to RabbitMQ", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		os.Getenv("NPM_MAHASISWA"),
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to open a channel", err)
		return
	}

	counter := 10
	for counter <= 100 {
		body := fmt.Sprintf("%d%%", counter)
		fmt.Printf("sending %s\n", body)
		err = ch.Publish(
			os.Getenv("NPM_MAHASISWA"), // name
			"dummyRouting",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		counter += 10
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	go rmqSend()
}

func setupRoutes() {
	http.HandleFunc("/", handler)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER2_HOST"), nil))
}
