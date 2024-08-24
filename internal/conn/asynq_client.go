package conn

import (
	"github.com/hibiken/asynq"
	"os"
)

func CreateAsyncQClient() (*asynq.Client, error) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: os.Getenv("ASYNCQ_REDIS_ADDR")})
	defer client.Close()
	return client, nil
}
