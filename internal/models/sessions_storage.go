package models

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/internal/conn"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/go-redis/redis"
	"log"
)

type SessionModel struct {
	Client *redis.Client
}

func (s *SessionModel) createCookie(user *payloads.UserRecord) map[string]string {
	return map[string]string{
		"username": user.UserName,
		"password": user.Password,
	}
}

func (s *SessionModel) createKey(userName string) string {
	return userName + "-session"
}

func (s *SessionModel) WriteSession(userName string, user *payloads.UserRecord) error {
	// Errors should be handled here
	key := s.createKey(userName)
	value := s.createCookie(user)

	bs, _ := json.Marshal(value)
	err := s.Client.Set(key, bs, conn.Ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionModel) ReadSession(userName string) (map[string]string, error) {
	key := s.createKey(userName)
	valueStr, err := s.Client.Get(key).Result()

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

func (s *SessionModel) DeleteSession(userName string) error {
	return s.Client.Del(userName).Err()
}
