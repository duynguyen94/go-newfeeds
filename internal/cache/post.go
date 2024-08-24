package cache

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/go-redis/redis"
	"strconv"
)

type PostCache interface {
	// WritePosts Set a list of user's post to cache
	WritePosts(userId int, posts []payloads.PostPayload) error

	// ReadPosts get a list of user's post
	ReadPosts(userId int) ([]payloads.PostPayload, error)
}

func NewPostCache(client *redis.Client) PostCache {
	return &postCache{client: client}
}

type postCache struct {
	client *redis.Client
}

func (p *postCache) createKey(userId int) string {
	return strconv.Itoa(userId) + "-newsfeed"
}

func (p *postCache) WritePosts(userId int, posts []payloads.PostPayload) error {
	key := p.createKey(userId)
	bs, err := json.Marshal(posts)

	if err != nil {
		return err
	}

	err = p.client.Set(key, bs, Ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (p *postCache) ReadPosts(userId int) ([]payloads.PostPayload, error) {
	key := p.createKey(userId)

	valueStr, err := p.client.Get(key).Result()

	// Empty
	if err == redis.Nil {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var res []payloads.PostPayload
	err = json.Unmarshal([]byte(valueStr), &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}
