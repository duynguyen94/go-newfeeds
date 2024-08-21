package main

import (
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/conn"
	"github.com/duynguyen94/go-newfeeds/internal/database"
	"github.com/duynguyen94/go-newfeeds/internal/newsfeed"
	"github.com/duynguyen94/go-newfeeds/internal/newsfeed/async"
	"github.com/duynguyen94/go-newfeeds/internal/object_store"
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

	_, err = object_store.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	asyncqClient, err := conn.CreateAsyncQClient()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	postCache := cache.NewPostCache(cacheClient)
	postDB := database.NewPostDB(db)
	userDB := database.NewUserDB(db)

	asyncTask := async.NewNewsfeedAsync(asyncqClient, postDB, userDB, postCache)
	newsfeedHandler := newsfeed.NewHandler(postCache, asyncTask)
	newsfeed.RouteV1(newsfeedHandler, r)

	r.Run(":8081")
}
