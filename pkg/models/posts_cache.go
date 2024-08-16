package models

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/pkg/conn"
	"github.com/go-redis/redis"
	"strconv"
)

type PostCacheModel struct {
	Client *redis.Client
}

func (p *PostCacheModel) createKey(userId int) string {
	return strconv.Itoa(userId) + "-newsfeed"
}

func (p *PostCacheModel) WritePost(userId int, posts []PostRecord) error {
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

func (p *PostCacheModel) ReadPost(userId int) ([]PostRecord, error) {
	key := p.createKey(userId)

	valueStr, err := p.Client.Get(key).Result()

	// Empty
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var res []PostRecord
	err = json.Unmarshal([]byte(valueStr), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
