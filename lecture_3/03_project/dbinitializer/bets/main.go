package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Note: due to the usage of relative paths, this script has to be run from this directory (go run main.go).
	// Running from Goland directly may cause incorrect behaviour.

	log.Println("If exists, dropping existing bets database...")
	os.Remove("../../db/bets.db")

	log.Println("Creating bets database...")
	file, err := os.Create("../../db/bets.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("bets database created")

	betsDatabase, _ := sql.Open("sqlite3", "../../db/bets.db")
	defer betsDatabase.Close()
	createBetsTable(betsDatabase)
}

func createBetsTable(db *sql.DB) {
	createBetsTableSQL := `CREATE TABLE bets (
		"id" TEXT NOT NULL PRIMARY KEY,
		"customer_id" TEXT NOT NULL,
		"status" TEXT NOT NULL,
		"selection_id" TEXT NOT NULL,
		"selection_coefficient" INTEGER NOT NULL,
		"payment" INTEGER NOT NULL,
		"payout" INTEGER
	  );`

	log.Println("Creating bets table...")
	statement, err := db.Prepare(createBetsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("bets table created")
}
