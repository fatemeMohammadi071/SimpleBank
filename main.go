package main

import (
	"database/sql"
	"log"

	"github.com/simpleBank/util"

	_ "github.com/lib/pq"
	"github.com/simpleBank/api"
	db "github.com/simpleBank/db/sqlc"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("con not load configuration", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("connection failed", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddres)

	if err != nil {
		log.Fatal("con not start server", err)
	}
}
