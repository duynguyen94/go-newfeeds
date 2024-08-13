package main

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"log"
	"time"
)

const (
	ttl       = time.Minute * 15
	redisHost = "localhost:6379"
	redisDb   = 0
)

func userRecordtoCookie(user *UserRecord) map[string]string {
	return map[string]string{
		"username": user.UserName,
		"password": user.Password,
	}
}

func createRedisClient() (*redis.Client, error) {
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

func SetJSONRedis(redisClient *redis.Client, key string, value *map[string]string) error {
	// Errors should be handled here
	bs, _ := json.Marshal(value)
	err := redisClient.Set(key, bs, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJSONRedis(redisClient *redis.Client, key string) (map[string]string, error) {
	valueStr, err := redisClient.Get(key).Result()

	// Empty
	if err == redis.Nil {
		log.Println("Key not found")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var res map[string]string
	err = json.Unmarshal([]byte(valueStr), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
