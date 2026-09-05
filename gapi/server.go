package gapi

import (
	"fmt"

	db "github.com/simpleBank/db/sqlc"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/token"
	"github.com/simpleBank/util"
	"github.com/simpleBank/worker"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config          util.Config
	tokenMaker      token.Maker
	store           db.Store
	taskDistributor worker.TaskDistributor
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		config:          config,
		store:           store,
		tokenMaker:      tokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
