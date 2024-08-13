package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	users UserDBModel
}

// SignUp
func (e *Env) SignUpHandler(c *gin.Context) {
	// Parse request body
	var newUserRecord UserRecord
	err := c.ShouldBindBodyWithJSON(&newUserRecord)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// Create new user
	err = newUserRecord.DOBtoDate()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	newUserRecord.HashPassword()
	id, err := e.users.CreateNewUser(newUserRecord)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}

// TODO Login
func (e *Env) LoginHandler(c *gin.Context) {
	// TODO Parse body
	// TODO Check redis if user existed
	// TODO check users exists + pass matched
	// TODO Gen cookie + Store in redis
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

// TODO Edit profile
func (e *Env) EditProfileHandler(c *gin.Context) {
	// Parse request body
	var newUserRecord UserRecord
	err := c.ShouldBindBodyWithJSON(&newUserRecord)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	user, err := e.users.GetUserRecord(userId)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	user.Merge(&newUserRecord)

	err = e.users.OverwriteUserRecord(userId, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// TODO Remove data from caching

	c.JSON(http.StatusOK, gin.H{
		"message": "updated successfully",
	})
}

func main() {
	// Setup shared connection,
	// follow https://www.alexedwards.net/blog/organising-database-access
	db, err := initDBConn()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		users: UserDBModel{DB: db},
	}

	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.POST("/v1/users/login", env.LoginHandler)
	r.POST("/v1/users", env.SignUpHandler)
	r.PUT("/v1/users/:id", env.EditProfileHandler)

	r.Run()
}
