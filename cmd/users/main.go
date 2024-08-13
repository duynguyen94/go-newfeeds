package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
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
func LoginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

// TODO Edit profile
func EditProfileHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
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

	r.POST("/v1/users/login", LoginHandler)
	r.POST("/v1/users", env.SignUpHandler)
	r.PUT("/v1/users", EditProfileHandler)

	r.Run()
}
