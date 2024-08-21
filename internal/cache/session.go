package cache

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/go-redis/redis"
	"log"
)

type SessionCache interface {
	// Write to cache
	Write(userName string, payload *payloads.UserPayload) error

	// Read from cache
	Read(userName string) (map[string]string, error)

	// Delete from cache
	Delete(userName string) error
}

func NewSessionCache(client *redis.Client) SessionCache {
	return &sessionCache{client: client}
}

type sessionCache struct {
	client *redis.Client
}

func (s *sessionCache) createCookie(user *payloads.UserPayload) map[string]string {
	return map[string]string{
		"username": user.UserName,
		"password": user.Password,
	}
}

func (s *sessionCache) createKey(userName string) string {
	return userName + "-session"
}

func (s *sessionCache) Write(userName string, payload *payloads.UserPayload) error {
	// Errors should be handled here
	key := s.createKey(userName)
	value := s.createCookie(payload)

	bs, _ := json.Marshal(value)
	err := s.client.Set(key, bs, Ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionCache) Read(userName string) (map[string]string, error) {
	key := s.createKey(userName)
	valueStr, err := s.client.Get(key).Result()

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

func (s *sessionCache) Delete(userName string) error {
	return s.client.Del(userName).Err()
}
