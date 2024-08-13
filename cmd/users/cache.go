package main

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/cmd/users/repo"
	"github.com/go-redis/redis"
	"log"
)

type SessionModel struct {
	cache *redis.Client
}

func (s *SessionModel) createCookie(user *UserRecord) map[string]string {
	return map[string]string{
		"username": user.UserName,
		"password": user.Password,
	}
}

func (s *SessionModel) createKey(userName string) string {
	return userName + "-session"
}

func (s *SessionModel) WriteSession(userName string, user *UserRecord) error {
	// Errors should be handled here
	key := s.createKey(userName)
	value := s.createCookie(user)

	bs, _ := json.Marshal(value)
	err := s.cache.Set(key, bs, repo.Ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionModel) ReadSession(userName string) (map[string]string, error) {
	key := s.createKey(userName)
	valueStr, err := s.cache.Get(key).Result()

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

func (s *SessionModel) deleteSession(userName string) error {
	return s.cache.Del(userName).Err()
}
