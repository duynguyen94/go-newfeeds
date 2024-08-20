package async

import (
	"encoding/json"
	"github.com/duynguyen94/go-newfeeds/pkg/models"
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

type TaskProcessor struct {
	Client *asynq.Client
	Users  models.UserDBModel
	Posts  models.PostCacheModel
}

func (processor *TaskProcessor) GenNewsfeedTasks(userId int) (*asynq.Task, error) {
	payload, err := json.Marshal(newsfeedTaskPayload{userId: userId})
	if err != nil {
		return nil, err
	}

	task := asynq.NewTask(TypeGenNewsfeed, payload, asynq.MaxRetry(5), asynq.Timeout(20*time.Minute))
	return task, nil
}

func (processor *TaskProcessor) GenNewsfeed(userId int) error {
	userPosts, err := processor.Users.ViewPosts(userId)
	if err != nil {
		log.Println(err)
		return err
	}

	friendPosts, err := processor.Users.ViewFriendPost(userId)
	if err != nil {
		log.Println(err)
		return err
	}

	userPosts = append(userPosts, friendPosts...)

	err = processor.Posts.WritePost(userId, userPosts)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("Finished generating newsfeed")
	return nil
}
