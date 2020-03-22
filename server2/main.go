package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

func rmqSend(file multipart.File, routingKey string) {
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

	fBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("Failed to read file:", err)
		return
	}

	var buff bytes.Buffer
	gzw, err := gzip.NewWriterLevel(&buff, gzip.BestCompression)
	if err != nil {
		log.Println("Failed to create gzip writer", err)
		return
	}
	defer gzw.Close()

	part := 1
	partLen := (len(fBytes) + 9) / 10

	for part <= 10 {
		left := (part-1) * partLen
		right := part * partLen

		if part == 10 {
			right = len(fBytes)
		}

		log.Printf("Compressing %d%%", part*10)

		_, err := gzw.Write(fBytes[left:right])
		if err != nil {
			log.Println("Failed to write part ", part, err)
			return
		}

		body := fmt.Sprintf("%d%%", part*10)

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

		part += 1
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	routingKey := r.Header.Get("X-ROUTING-KEY")

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("Failed to get the file:", err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	go rmqSend(file, routingKey)
}

func setupRoutes() {
	http.HandleFunc("/", handler)
}

func main() {
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER2_HOST"), nil))
}
