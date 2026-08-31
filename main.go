package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/simpleBank/gapi"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	runGRPCServer(config, store)
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddres)
	if err != nil {
		log.Fatal("Can not create listener", err)
	}

	log.Printf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("Can not start grpc server", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal("Can not create server: ", err)
	}
	err = server.Start(config.HTTPServerAddres)

	if err != nil {
		log.Fatal("Can not start server", err)
	}
}
