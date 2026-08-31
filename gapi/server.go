package gapi

import (
	"fmt"

	db "github.com/simpleBank/db/sqlc"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/token"
	"github.com/simpleBank/util"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	tokenMaker token.Maker
	store      db.Store
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
