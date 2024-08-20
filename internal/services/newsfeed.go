package services

import (
	"github.com/duynguyen94/go-newfeeds/internal/async"
	models2 "github.com/duynguyen94/go-newfeeds/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type NewsfeedServices struct {
	Users        models2.UserDBModel
	Posts        models2.PostCacheModel
	ImageStorage models2.ImagePostStorageModel
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
