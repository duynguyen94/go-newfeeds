package main

import (
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/database"
	"github.com/duynguyen94/go-newfeeds/internal/object_store"
	"github.com/duynguyen94/go-newfeeds/internal/post"
	"github.com/duynguyen94/go-newfeeds/internal/user"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting user services")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Setup shared connection,
	// follow https://www.alexedwards.net/blog/organising-database-access
	db, err := database.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := cache.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	minIOClient, err := object_store.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20 // Max 8MB

	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	//postCache := cache.NewPostCache(cacheClient)
	sessionCache := cache.NewSessionCache(cacheClient)
	postDB := database.NewPostDB(db)
	userDB := database.NewUserDB(db)
	imageStorage := object_store.NewImageStorage(minIOClient, "images")

	// Simple ping
	err = imageStorage.BucketExists()
	if err != nil {
		log.Fatal(err)
	}

	userHandler := user.NewHandler(userDB, sessionCache)
	user.RouteV1(userHandler, r)

	postHandler := post.NewHandler(postDB, imageStorage)
	post.RouteV1(postHandler, r)

	// Performance stats
	pprof.Register(r)
	r.Run()
}
