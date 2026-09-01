package main

import (
	"context"
	"database/sql"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/simpleBank/doc"
	"github.com/simpleBank/gapi"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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

	// usin this line defin two routes for grpc and http that dose not block each other
	go runGatewayServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Can not create server", err)
	}

	jsonOptions := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOptions)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot create register")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerFS, err := fs.Sub(doc.SwaggerFiles, "swagger")
	if err != nil {
		log.Fatal("cannot create swagger sub filesystem", err)
	}
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServerFS(swaggerFS)))

	listener, err := net.Listen("tcp", config.HTTPServerAddres)
	if err != nil {
		log.Fatal("Can not create listener", err)
	}

	log.Printf("start http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("Can not start http gateway server", err)
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
