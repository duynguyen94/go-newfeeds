package main

import (
	"github.com/duynguyen94/go-newfeeds/internal/async"
	"github.com/duynguyen94/go-newfeeds/internal/conn"
	models2 "github.com/duynguyen94/go-newfeeds/internal/models"
	"github.com/duynguyen94/go-newfeeds/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting newsfeed services")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := conn.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := conn.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	minIOClient, err := conn.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	asyncqClient, err := conn.CreateAsyncQClient()
	if err != nil {
		log.Fatal(err)
	}

	app := &services.NewsfeedServices{
		Users:        models2.UserDBModel{DB: db},
		Posts:        models2.PostCacheModel{Client: cacheClient},
		ImageStorage: models2.ImagePostStorageModel{Client: minIOClient, Bucket: models2.DefaultBucket},
		Tasks: async.TaskProcessor{
			Client: asyncqClient,
			Users:  models2.UserDBModel{DB: db},
			Posts:  models2.PostCacheModel{Client: cacheClient},
		},
	}

	// Simple ping
	err = app.ImageStorage.BucketExists()
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

	r.Run(":8081")
}
