package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// TODO Login
	r.POST("/v1/users/login", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Need implementation",
		})
	})

	// TODO SignUp
	r.POST("/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Need implementation",
		})
	})

	// TODO Edit profile
	r.PUT("/v1/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Need implementation",
		})
	})

	r.Run()
}
