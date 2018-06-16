package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/hoop33/roster/players"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	db, err := createDatabase()
	if err != nil {
		log.Println("failed to connect to database:", err)
		os.Exit(1)
	}
	defer db.Close()

	ps := createPlayersService(db)

	players, err := ps.ListPlayers(context.Background(), "QB")
	if err != nil {
		log.Println("failed to list players:", err)
		os.Exit(1)
	}

	for _, player := range players {
		log.Println(player.String())
	}
}

func createDatabase() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=roster sslmode=disable",
		os.Getenv("ROSTER_USER"),
		os.Getenv("ROSTER_PASSWORD")))
	if err != nil {
		return nil, err
	}

	ddl, err := ioutil.ReadFile("create_table.sql")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(ddl))
	return db, err
}

func createPlayersService(db *sqlx.DB) players.Service {
	ps := players.NewService(db)
	return ps
}
