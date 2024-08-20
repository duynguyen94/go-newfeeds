package models

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/internal/conn"
	models2 "github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/go-redis/redis"
	"strconv"
)

type PostCacheModel struct {
	Client *redis.Client
}

func (p *PostCacheModel) createKey(userId int) string {
	return strconv.Itoa(userId) + "-newsfeed"
}

func (p *PostCacheModel) WritePost(userId int, posts []models2.PostRecord) error {
	key := p.createKey(userId)
	bs, err := json.Marshal(posts)

	if err != nil {
		return err
	}

	err = p.Client.Set(key, bs, conn.Ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *PostCacheModel) ReadPost(userId int) ([]models2.PostRecord, error) {
	key := p.createKey(userId)

	valueStr, err := p.Client.Get(key).Result()

	// Empty
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var res []models2.PostRecord
	err = json.Unmarshal([]byte(valueStr), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
