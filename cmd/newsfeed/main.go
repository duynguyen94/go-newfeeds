package main

import (
	repo2 "github.com/duynguyen94/go-newfeeds/pkg/conn"
	"github.com/duynguyen94/go-newfeeds/pkg/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	users        models.UserDBModel
	posts        models.PostCacheModel
	imageStorage models.ImagePostStorageModel
	tasks        TaskProcessor
}

func (app *App) GetNewsfeedsHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	posts, err := app.posts.ReadPost(userId)

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

func (app *App) GenNewsfeedHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = app.tasks.GenNewsfeed(userId)
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

func main() {
	db, err := repo2.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := repo2.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	minIOClient, err := repo2.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	asyncqClient, err := repo2.CreateAsyncQClient()
	if err != nil {
		log.Fatal(err)
	}

	app := &App{
		users:        models.UserDBModel{DB: db},
		posts:        models.PostCacheModel{Client: cacheClient},
		imageStorage: models.ImagePostStorageModel{Client: minIOClient, Bucket: models.DefaultBucket},
		tasks: TaskProcessor{
			client: asyncqClient,
			users:  models.UserDBModel{DB: db},
			posts:  models.PostCacheModel{Client: cacheClient},
		},
	}

	// Simple ping
	err = app.imageStorage.BucketExists()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	r.GET("/v1/newsfeeds/:id", app.GetNewsfeedsHandler)
	r.POST("/v1/newsfeeds/:id", app.GenNewsfeedHandler)

	r.Run()
}
