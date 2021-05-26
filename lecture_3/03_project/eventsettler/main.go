package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"

	_ "github.com/mattn/go-sqlite3"
)

const numberOfSelectionsToSettle = 10

type eventUpdateDto struct {
	Id      string `json:"id"`
	Outcome string `json:"outcome"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Note: this script was written in a way "just make it work", it cannot be
	// considered as a well-written code.

	// Note: due to the usage of relative paths, this script has to be run from this directory (go run main.go).
	// Running from Goland directly may cause incorrect behaviour.

	rand.Seed(time.Now().UnixNano())

	// Set up RabbitMQ queue.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"event-updates", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	failOnError(err, "failed to declare a queue")

	// Load the desired number of selections which are still active.
	selectionIds := make([]string, 0)
	sqliteDatabase, err := sql.Open("sqlite3", "../db/bets.db")
	failOnError(err, "failed to open connection to DB")
	defer sqliteDatabase.Close() // Defer Closing the database
	rows, err := sqliteDatabase.Query(
		`SELECT DISTINCT selection_id FROM bets WHERE status='active' LIMIT ` + strconv.Itoa(numberOfSelectionsToSettle) + ";",
	)
	failOnError(err, "failed to execute a SQL statement")
	defer rows.Close()
	for rows.Next() {
		var id string
		rows.Scan(&id)
		selectionIds = append(selectionIds, id)
	}

	for _, selectionId := range selectionIds {
		var outcome string
		if rand.Float64() > 0.5 {
			outcome = "lost"
		} else {
			outcome = "won"
		}

		eventUpdate := &eventUpdateDto{
			Id: selectionId,
			Outcome: outcome,
		}

		eventUpdateJson, err := json.Marshal(eventUpdate)
		failOnError(err, "failed to marshal an event update")

		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        eventUpdateJson,
			},
		)
		failOnError(err, "failed to publish a message")
		log.Printf("Sent %s", eventUpdateJson)
	}
}
