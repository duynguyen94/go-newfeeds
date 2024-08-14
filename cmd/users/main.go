package main

import (
	"errors"
	"github.com/duynguyen94/go-newfeeds/cmd/users/repo"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Env struct {
	users      UserDBModel
	sessions   SessionModel
	imageStore ImageStorageModel
}

// TODO Find the way to do it properly
type friendPayload struct {
	FriendId int `json:"friendId"`
}

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

func (e *Env) LoginHandler(c *gin.Context) {
	var user UserRecord
	err := c.ShouldBindBodyWithJSON(&user)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	sess, err := e.sessions.ReadSession(user.UserName)
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
	curUser, err := e.users.GetUserRecordByUsername(user.UserName)

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

	err = e.sessions.WriteSession(user.UserName, &user)
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

	err = e.sessions.deleteSession(user.UserName)
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

func (e *Env) FollowHandler(c *gin.Context) {
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

	err = e.users.FollowUser(userId, requestBody.FriendId)
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

func (e *Env) UnFollowHandler(c *gin.Context) {
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

	err = e.users.UnFollowUser(userId, requestBody.FriendId)
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

func (e *Env) GetFollowers(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	followers, err := e.users.ViewFollowers(userId)
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

func (e *Env) ViewFriendPostsHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	posts, err := e.users.ViewFriendPost(userId)

	// TODO For each posts, add gen pre-signed url for image
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func (e *Env) GetPost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func (e *Env) CreatePost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}
func (e *Env) EditPost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}
func (e *Env) DeletePost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func (e *Env) LikePost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func (e *Env) CommentPost(c *gin.Context) {
	// TODO
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func main() {
	// Setup shared connection,
	// follow https://www.alexedwards.net/blog/organising-database-access
	db, err := repo.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := repo.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		users:    UserDBModel{DB: db},
		sessions: SessionModel{cache: cacheClient},
	}

	r := gin.Default()
	r.GET("/health-check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// TODO Organize API with group

	// Users
	r.POST("/v1/users/login", env.LoginHandler)
	r.POST("/v1/users", env.SignUpHandler)
	r.PUT("/v1/users/:id", env.EditProfileHandler)

	// Friends
	r.GET("/v1/friends/:id", env.GetFollowers)
	r.POST("/v1/friends/:id", env.FollowHandler)
	r.DELETE("/v1/friends/:id", env.UnFollowHandler)
	r.GET("/v1/friends/:id/posts", env.ViewFriendPostsHandler)

	// Posts
	r.GET("/v1/posts/:post_id", env.GetPost)
	r.POST("/v1/posts", env.CreatePost)
	r.PUT("/v1/posts/:post_id", env.EditPost)
	r.DELETE("/v1/posts/:post_id", env.DeletePost)
	r.POST("/v1/posts/:post_id/comments", env.CommentPost)
	r.POST("/v1/posts/:post_id/likes", env.LikePost)

	r.Run()
}
