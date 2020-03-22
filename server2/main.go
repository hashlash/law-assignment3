package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
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
	conn, err := amqp.Dial(os.Getenv("AMQP_URL"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to declare an exchange")

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

func main() {
	rmqSend()
}
