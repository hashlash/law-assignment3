package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"time"
)

func rmqSend(routingKey string) {
	time.Sleep(3 * time.Second)
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	if err != nil {
		log.Println("Failed to connect to RabbitMQ:", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println("Failed to open a channel:", err)
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
		log.Println("Failed to open a channel:", err)
		return
	}

	counter := 10
	for counter <= 100 {
		body := fmt.Sprintf("%d%%", counter)
		fmt.Printf("sending %s\n", body)

		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		}

		err = ch.Publish(
			os.Getenv("NPM_MAHASISWA"), // name
			routingKey,
			false,
			false,
			msg,
		)
		if err != nil {
			log.Println("Failed to publish:", err)
			continue
		}
		counter += 10
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	routingKey := r.Header.Get("X-ROUTING-KEY")

	go rmqSend(routingKey)
}

func setupRoutes() {
	http.HandleFunc("/", handler)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER2_HOST"), nil))
}
