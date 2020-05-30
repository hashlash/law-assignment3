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
)

func rmqSend(file multipart.File, routingKey string) {
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

	rCh := make(chan amqp.Return)
	ch.NotifyReturn(rCh)

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
			part == 10,
			false,
			msg,
		)
		if err != nil {
			log.Println("Failed to publish:", err)
			continue
		}

		part += 1
	}

	for r := range rCh {
		log.Println("Failed to publish:", r)
		log.Println("Retrying")

		msg := amqp.Publishing{
			Headers:         r.Headers,
			ContentType:     r.ContentType,
			ContentEncoding: r.ContentEncoding,
			DeliveryMode:    r.DeliveryMode,
			Priority:        r.Priority,
			CorrelationId:   r.CorrelationId,
			ReplyTo:         r.ReplyTo,
			Expiration:      r.Expiration,
			MessageId:       r.MessageId,
			Timestamp:       r.Timestamp,
			Type:            r.Type,
			UserId:          r.UserId,
			AppId:           r.AppId,
			Body:            r.Body,
		}

		err = ch.Publish(
			r.Exchange,
			r.RoutingKey,
			true,
			false,
			msg,
		)
		if err != nil {
			log.Println("Failed to publish:", err)
		}
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
	log.Println("Serving on:", os.Getenv("SERVER2_HOST"))
	setupRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER2_HOST"), nil))
}
