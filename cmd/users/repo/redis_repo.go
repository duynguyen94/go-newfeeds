package repo

import (
	"github.com/go-redis/redis"
	"time"
)

const (
	Ttl       = time.Minute * 15
	redisHost = "localhost:6379"
	redisDb   = 0
)

func CreateRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
		DB:   redisDb,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
