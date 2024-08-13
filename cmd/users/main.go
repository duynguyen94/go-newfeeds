package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// TODO Login
func LoginHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

// TODO SignUp
func SignUpHandler(c *gin.Context) {
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
	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	r.Group("v1")
	{
		r.POST("/users/login", LoginHandler)
		r.POST("/users", SignUpHandler)
		r.PUT("/users", EditProfileHandler)
	}

	r.Run()
}
