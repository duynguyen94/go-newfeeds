package main

import (
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/conn"
	"github.com/duynguyen94/go-newfeeds/internal/database"
	models2 "github.com/duynguyen94/go-newfeeds/internal/models"
	"github.com/duynguyen94/go-newfeeds/internal/newsfeed"
	"github.com/duynguyen94/go-newfeeds/internal/newsfeed/async"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	log.Println("Starting newsfeed services")
	// Refactor with example from https://github.com/zacscoding/gin-rest-api-example/blob/master/internal/article/database/article.go
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	_, err = conn.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	asyncqClient, err := conn.CreateAsyncQClient()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	postCache := cache.PostCacheModel{Client: cacheClient}
	asyncTask := async.TaskProcessor{
		Client: asyncqClient,
		Users:  models2.UserDBModel{DB: db},
		Posts:  cache.PostCacheModel{Client: cacheClient},
	}
	newsfeedHandler := newsfeed.NewHandler(postCache, asyncTask)
	newsfeed.RouteV1(newsfeedHandler, r)

	r.Run(":8081")
}
