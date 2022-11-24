package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Request struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

type Response struct {
	Path []string `json:"path"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func findPathRPC(source, destination string) (res []string, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := Request{source, destination}
	req_bytes, err := json.Marshal(req)
	failOnError(err, "Failed to marshal request")

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          req_bytes,
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			var resp Response
			err := json.Unmarshal(d.Body, &resp)
			failOnError(err, "Failed to unmarshal response")

			res = resp.Path
			break
		}
	}

	return
}

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: client [source] [destination]")
	}

	path, _ := findPathRPC(os.Args[1], os.Args[2])

	for i := len(path) - 1; i >= len(path)/2; i-- {
		path[i], path[len(path)-i-1] = path[len(path)-i-1], path[i]

	}

	fmt.Println(strings.Join(path, " --> "))
}
