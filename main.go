package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/simpleBank/api"
	db "github.com/simpleBank/db/sqlc"
)

const (
	dbDriver     = "postgres"
	dbSource     = "postgres://root:Fateme025@localhost:5432/simple_bank?sslmode=disable"
	serverAddres = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("connection failed")
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddres)

	if err != nil {
		log.Fatal("con not start server", err)
	}
}
