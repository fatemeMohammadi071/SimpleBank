package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

// TaskDistributor pushes background tasks onto the Redis queue.
// The API/gRPC server depends on this interface so it can be mocked in tests.
type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(
		ctx context.Context,
		payload *PayloadSendVerifyEmail,
		opts ...asynq.Option,
	) error
}

// RedisTaskDistributor is the production implementation backed by asynq.
type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	return &RedisTaskDistributor{
		client: client,
	}
}
