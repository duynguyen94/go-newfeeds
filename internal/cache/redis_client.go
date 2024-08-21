package cache

import (
	"github.com/go-redis/redis"
	"os"
	"strconv"
	"time"
)

const (
	Ttl = time.Minute * 15
)

func CreateRedisClient() (*redis.Client, error) {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		DB:   db,
	})

	_, err = client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return client, nil
}
