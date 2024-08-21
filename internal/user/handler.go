package user

import (
	"errors"
	"github.com/duynguyen94/go-newfeeds/internal/cache"
	"github.com/duynguyen94/go-newfeeds/internal/database"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	userDB       database.UserDB
	sessionCache cache.SessionCache
}

type friendPayload struct {
	FriendId int `json:"friendId"`
}

func (e *Handler) SignUp(c *gin.Context) {
	// Parse request body
	var newUserRecord payloads.UserPayload
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
	id, err := e.userDB.New(&newUserRecord)

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

func (e *Handler) Login(c *gin.Context) {
	var user payloads.UserPayload
	err := c.ShouldBindBodyWithJSON(&user)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	sess, err := e.sessionCache.Read(user.UserName)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	if sess != nil {
		// Check pass
		if sess["password"] != user.Password {
			err = errors.New("wrong password")
			log.Println(err.Error())
			c.JSON(http.StatusForbidden, gin.H{
				"message": "err: " + err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "Login successfully",
			})
		}
		return
	}

	// Login check
	curUser, err := e.userDB.GetByUsername(user.UserName)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	if curUser.IsMatchPassword(user.Password) == false {
		err = errors.New("wrong password")
		log.Println(err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.sessionCache.Write(user.UserName, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusForbidden, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successfully",
	})
}

func (e *Handler) EditProfile(c *gin.Context) {
	// Parse request body
	var newUserRecord payloads.UserPayload
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

	user, err := e.userDB.GetById(userId)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	user.Merge(&newUserRecord)

	err = e.userDB.Edit(userId, &user)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.sessionCache.Delete(user.UserName)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated successfully",
	})
}

func (e *Handler) Follow(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var requestBody friendPayload
	err = c.ShouldBindBodyWithJSON(&requestBody)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.userDB.Follow(userId, requestBody.FriendId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Follow successfully",
	})
}

func (e *Handler) UnFollow(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var requestBody friendPayload
	err = c.ShouldBindBodyWithJSON(&requestBody)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.userDB.UnFollow(userId, requestBody.FriendId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "UnFollow successfully",
	})
}

func (e *Handler) GetFollowers(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	followers, err := e.userDB.GetFollowers(userId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
	})
}

// RouteV1 Route for post service
func RouteV1(h *Handler, r *gin.Engine) {
	v1 := r.Group("v1")
	v1.Use()
	{
		v1.POST("/users/login", h.Login)
		v1.POST("/users", h.SignUp)
		v1.PUT("/users/:id", h.EditProfile)

		v1.GET("/friends/:id", h.GetFollowers)
		v1.POST("/friends/:id", h.Follow)
		v1.DELETE("/friends/:id", h.UnFollow)

		// TODO
		//v1.GET("/friends/:id/posts", h.ViewFriendPostsHandler)
	}
}

func NewHandler(userDB database.UserDB, sessionCache cache.SessionCache) *Handler {
	return &Handler{
		userDB:       userDB,
		sessionCache: sessionCache,
	}
}
