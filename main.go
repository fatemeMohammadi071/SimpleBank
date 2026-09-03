package main

import (
	"context"
	"database/sql"
	"io/fs"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/simpleBank/doc"
	"github.com/simpleBank/gapi"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/simpleBank/api"
	migration "github.com/simpleBank/db/migration"
	db "github.com/simpleBank/db/sqlc"
)

func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal().Err(err).Msg("con not load configuration")
	}

	if config.Enviroment == "Development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal().Err(err).Msg("connection failed")
	}

	runDBMigration(conn)

	store := db.NewStore(conn)

	// usin this line defin two routes for grpc and http that dose not block each other
	go runGatewayServer(config, store)
	runGRPCServer(config, store)
}

func runDBMigration(conn *sql.DB) {
	source, err := iofs.New(migration.FS, ".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration source")
	}

	driver, err := postgres.WithInstance(conn, &postgres.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migration driver")
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create migrate instance")
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}

func runGRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddres)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create listener")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not start grpc server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create server")
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
		log.Fatal().Err(err).Msg("cannot create register")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	swaggerFS, err := fs.Sub(doc.SwaggerFiles, "swagger")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create swagger sub filesystem")
	}
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServerFS(swaggerFS)))

	listener, err := net.Listen("tcp", config.HTTPServerAddres)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not create listener")
	}

	handler := gapi.HttpLogger(mux)

	log.Info().Msgf("start http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Can not start http gateway server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		log.Fatal().Err(err).Msg("Can not create server")
	}
	err = server.Start(config.HTTPServerAddres)

	if err != nil {
		log.Fatal().Err(err).Msg("Can not start server")
	}
}
