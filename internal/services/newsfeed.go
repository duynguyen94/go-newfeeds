package services

import (
	"github.com/duynguyen94/go-newfeeds/internal/async"
	"github.com/duynguyen94/go-newfeeds/pkg/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type NewsfeedServices struct {
	Users        models.UserDBModel
	Posts        models.PostCacheModel
	ImageStorage models.ImagePostStorageModel
	Tasks        async.TaskProcessor
}

func (app *NewsfeedServices) GetNewsfeedsHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	posts, err := app.Posts.ReadPost(userId)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// TODO Gen downloadable image url
	// Client might retries several time before newsfeed appear
	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
	return
}

func (app *NewsfeedServices) GenNewsfeedHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = app.Tasks.GenNewsfeed(userId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Queued task",
	})

}
