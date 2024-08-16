package conn

import (
	"github.com/hibiken/asynq"
)

// TODO Load from env
const (
	redisAddr = "localhost:6379"
)

func CreateAsyncQClient() (*asynq.Client, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer client.Close()
	return client, nil
}
