package gapi

import (
	"context"
	"time"

	"github.com/hibiken/asynq"
	"github.com/lib/pq"
	db "github.com/simpleBank/db/sqlc"
	"github.com/simpleBank/pb"
	"github.com/simpleBank/util"
	"github.com/simpleBank/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violations := validateCreateUserRequest(req); violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)

	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username already exist: %s", err)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
	}

	// Send the verification email in the background so the client is not
	// blocked waiting for the mail server.
	taskPayload := &worker.PayloadSendVerifyEmail{
		Username: user.Username,
	}
	opts := []asynq.Option{
		asynq.MaxRetry(10),
		asynq.ProcessIn(10 * time.Second),
		asynq.Queue(worker.QueueCritical),
	}
	if err := server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to distribute task to send verify email: %s", err)
	}

	res := &pb.CreateUserResponse{
		User: convertUser(user),
	}

	return res, nil
}
