package main

import (
	"errors"
	repo2 "github.com/duynguyen94/go-newfeeds/pkg/conn"
	"github.com/duynguyen94/go-newfeeds/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Env struct {
	users    models.UserDBModel
	posts    models.PostDBModel
	sessions models.SessionModel
	images   models.ImagePostStorageModel
}

// TODO Find the way to do it properly
type friendPayload struct {
	FriendId int `json:"friendId"`
}

func (e *Env) SignUpHandler(c *gin.Context) {
	// Parse request body
	var newUserRecord models.UserRecord
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
	var user models.UserRecord
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
	var newUserRecord models.UserRecord
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

	err = e.sessions.DeleteSession(user.UserName)
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
	postId, err := strconv.Atoi(c.Param("post_id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	p, err := e.posts.GetPostById(postId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// Expiration, TODO Move it into somewhere
	expiration := time.Minute * 30
	err = p.GenSignedUrl(e.images, expiration)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post": p,
	})

}

func (e *Env) CreatePost(c *gin.Context) {
	var post models.PostRecord
	err := c.ShouldBindBodyWithJSON(&post)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	postId, err := e.posts.CreatePost(&post)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}
	defer triggerGenNewsfeed(post.UserId)

	c.JSON(http.StatusOK, gin.H{
		"postId": postId,
	})
}

func (e *Env) UploadImage(c *gin.Context) {
	// FIXME should merge 2 api create post into 1
	postId, err := strconv.Atoi(c.Param("post_id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// single file
	file, header, err := c.Request.FormFile("filename")
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	filename := header.Filename
	//fSize := header.Size

	out, err := os.Create("./tmp/" + filename)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	defer out.Close()
	_, err = io.Copy(out, file)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	fStat, err := out.Stat()
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// Move file cursor to the start, follow https://www.reddit.com/r/golang/comments/wv7hky/i_cant_upload_the_file_to_the_minio_bucket/
	_, err = out.Seek(0, io.SeekStart)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	log.Printf("Stats %d\n", fStat.Size())
	imagePath, err := e.images.PutImage(out, postId, filename, fStat.Size())
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.posts.UpdateImagePath(postId, imagePath)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"postId":    postId,
		"imagePath": imagePath,
	})
}

func (e *Env) EditPost(c *gin.Context) {
	var newPost models.PostRecord
	err := c.ShouldBindBodyWithJSON(&newPost)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	curPost, err := e.posts.GetPostById(postId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	curPost.Merge(&newPost)
	err = e.posts.OverwritePost(postId, &curPost)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	defer triggerGenNewsfeed(curPost.UserId)

	// TODO Trigger delete cache if possible
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func (e *Env) DeletePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.posts.DeletePost(postId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Delete post successfully",
	})
}

func (e *Env) LikePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var p models.PostRecord
	err = c.ShouldBindBodyWithJSON(&p)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = e.posts.LikePost(postId, p.UserId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Like post successfully",
	})
}

func (e *Env) CommentPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var cmt models.CommentRecord
	err = c.ShouldBindBodyWithJSON(&cmt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	cmtId, err := e.posts.CommentPost(postId, cmt.UserId, cmt.ContentText)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"commentId": cmtId,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Setup shared connection,
	// follow https://www.alexedwards.net/blog/organising-database-access
	db, err := repo2.InitMySQLDBConn()
	if err != nil {
		log.Fatal(err)
	}

	cacheClient, err := repo2.CreateRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	minIOClient, err := repo2.CreateMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		users:    models.UserDBModel{DB: db},
		posts:    models.PostDBModel{DB: db},
		sessions: models.SessionModel{Client: cacheClient},
		images:   models.ImagePostStorageModel{Client: minIOClient, Bucket: models.DefaultBucket},
	}

	// Simple ping
	err = env.images.BucketExists()
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
	r.POST("/v1/posts/:post_id/images", env.UploadImage)
	r.PUT("/v1/posts/:post_id", env.EditPost)
	r.DELETE("/v1/posts/:post_id", env.DeletePost)
	r.POST("/v1/posts/:post_id/comments", env.CommentPost)
	r.POST("/v1/posts/:post_id/likes", env.LikePost)

	r.Run()
}
