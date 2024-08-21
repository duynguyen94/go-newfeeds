package async

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/database"
	"github.com/hibiken/asynq"
	"log"
	"time"
)

const (
	TypeGenNewsfeed = "newsfeed:gen"
)

type newsfeedTaskPayload struct {
	userId int
}

type NewsfeedAsync interface {

	// CreateTask convert userId into task payload
	CreateTask(userId int) (*asynq.Task, error)

	// Generate newfeed for given userId
	Generate(userId int) error
}

func NewNewsfeedAsync(client *asynq.Client, postDB database.PostDB, userDB database.UserDB, postCache cache.PostCache) NewsfeedAsync {
	return &taskProcessor{client: client, postDB: postDB, userDB: userDB, postCache: postCache}
}

type taskProcessor struct {
	client    *asynq.Client
	postDB    database.PostDB
	userDB    database.UserDB
	postCache cache.PostCache
}

func (processor *taskProcessor) CreateTask(userId int) (*asynq.Task, error) {
	payload, err := json.Marshal(newsfeedTaskPayload{userId: userId})
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeGenNewsfeed, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute))
	return task, nil
}

func (processor *taskProcessor) Generate(userId int) error {
	userPosts, err := processor.postDB.ListPostByUserId(userId)
	if err != nil {
		log.Println(err)
		return err
	}

	// TODO list post by follower
	//friendPosts, err := processor.Users.ViewFriendPost(userId)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//userPosts = append(userPosts, friendPosts...)

	err = processor.postCache.WritePosts(userId, userPosts)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Finished generating newsfeed")
	return nil
}
