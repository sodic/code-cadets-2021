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

	log.Println("If exists, dropping existing calc_bets database...")
	os.Remove("../../db/calc_bets.db")

	log.Println("Creating calc_bets database...")
	file, err := os.Create("../../db/calc_bets.db")
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("calc_bets database created")

	calcBetsDatabase, _ := sql.Open("sqlite3", "../../db/calc_bets.db")
	defer calcBetsDatabase.Close()
	createBetsTable(calcBetsDatabase)
	createBetsIndexOnSelection(calcBetsDatabase)
}

func createBetsTable(db *sql.DB) {
	createBetsTableSQL := `CREATE TABLE bets (
		"id" TEXT NOT NULL PRIMARY KEY,
		"selection_id" TEXT NOT NULL,
		"selection_coefficient" TEXT NOT NULL,
		"payment" INTEGER NOT NULL
	  );`

	log.Println("Creating bets table...")
	statement, err := db.Prepare(createBetsTableSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("bets table created")
}

func createBetsIndexOnSelection(db *sql.DB) {
	createBetsIndexOnSelectionSQL := `CREATE INDEX selection_idx ON bets (selection_id);`

	log.Println("Creating selection_idx on bets table...")
	statement, err := db.Prepare(createBetsIndexOnSelectionSQL)
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec()
	log.Println("selection_idx on bets table created")
}