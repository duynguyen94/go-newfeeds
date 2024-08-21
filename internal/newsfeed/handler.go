package newsfeed

import (
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/newsfeed/async"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type NewsfeedHandler struct {
	postCache cache.PostCache
	asyncTask async.TaskProcessor
}

func (app *NewsfeedHandler) GetNewsfeeds(c *gin.Context) {
	// TODO Get userId from cookies
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	posts, err := app.postCache.ReadPosts(userId)

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

func (app *NewsfeedHandler) GenNewsfeed(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = app.asyncTask.GenNewsfeed(userId)
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

// RouteV1 Route for newsfeed service
func RouteV1(h *NewsfeedHandler, r *gin.Engine) {
	v1 := r.Group("v1")
	v1.Use()
	{
		v1.GET("/newsfeeds/:id", h.GetNewsfeeds)
		v1.POST("/newsfeeds/:id", h.GenNewsfeed)
	}
}

func NewHandler(postCache cache.PostCache, taskProcessor async.TaskProcessor) *NewsfeedHandler {
	return &NewsfeedHandler{
		postCache: postCache,
		asyncTask: taskProcessor,
	}
}
