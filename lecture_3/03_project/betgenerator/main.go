package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
)

const numberOfBets = 1000
const numberOfDistinctCustomers = 50
const numberOfDistinctSelections = 10

type selection struct {
	id          string
	coefficient float64
}

type betDto struct {
	Id                   string  `json:"id"`
	CustomerId           string  `json:"customerId"`
	SelectionId          string  `json:"selectionId"`
	SelectionCoefficient float64 `json:"selectionCoefficient"`
	Payment              float64 `json:"payment"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func getRandomUUID() string {
	id, err := uuid.NewV4()
	failOnError(err, "failed to generate a uuid")
	return id.String()
}

func main() {
	// Note: this script was written in a way "just make it work", it cannot be
	// considered as a well-written code.

	rand.Seed(time.Now().UnixNano())

	// Set up RabbitMQ queue.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"bets-received", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "failed to declare a queue")

	// Generate customer data.
	customers := make([]string, 0)
	for i := 0; i < numberOfDistinctCustomers; i++ {
		customerId := getRandomUUID()
		customers = append(customers, customerId)
	}

	// Generate selections data.
	selections := make([]selection, 0)
	for i := 0; i < numberOfDistinctSelections; i++ {
		selectionId := getRandomUUID()
		selectionCoefficient := 1 + rand.Float64()*(10-1) // arbitraty number [1, 10] https://stackoverflow.com/a/49747128
		selections = append(selections, selection{id: selectionId, coefficient: selectionCoefficient})
	}

	// Generate and publish bets.
	for i := 0; i < numberOfBets; i++ {
		betId := getRandomUUID()
		customerId := customers[rand.Intn(numberOfDistinctCustomers)]
		selection := selections[rand.Intn(numberOfDistinctSelections)]
		payment := 10 + rand.Float64()*(500-10) // arbitraty number [10, 500] https://stackoverflow.com/a/49747128

		bet := &betDto{
			Id:                   betId,
			CustomerId:           customerId,
			SelectionId:          selection.id,
			SelectionCoefficient: float64(int(selection.coefficient*100)) / 100, // we want 2 decimal places; https://stackoverflow.com/a/18391072
			Payment:              float64(int(payment*100)) / 100,               // we want 2 decimal places; https://stackoverflow.com/a/18391072
		}

		betJson, err := json.Marshal(bet)
		failOnError(err, "failed to marshal a bet")

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        betJson,
			},
		)
		failOnError(err, "failed to publish a message")
		log.Printf("Sent %s", betJson)
	}
}
