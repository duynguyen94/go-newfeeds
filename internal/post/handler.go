package post

import (
	"github.com/duynguyen94/go-newfeeds/internal/database"
	"github.com/duynguyen94/go-newfeeds/internal/object_store"
	"github.com/duynguyen94/go-newfeeds/internal/payloads"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Handler struct {
	postDB       database.PostDB
	imageStorage object_store.ImageStorage
}

type friendPayload struct {
	FriendId int `json:"friendId"`
}

func (app *Handler) ViewFriendPostsHandler(c *gin.Context) {
	userId, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	posts, err := app.postDB.ListPostByUserId(userId)

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

func (app *Handler) GetPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	p, err := app.postDB.GetPostById(postId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// Expiration, TODO Move it into somewhere
	//expiration := time.Minute * 30
	//err = p.imageStorage(p.ContentImagePath, expiration)
	//if err != nil {
	//	log.Println(err.Error())
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"message": "err: " + err.Error(),
	//	})
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"post": p,
	})

}

func (app *Handler) CreatePost(c *gin.Context) {
	var post payloads.PostPayload
	err := c.ShouldBindBodyWithJSON(&post)

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	postId, err := app.postDB.New(&post)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}
	// TODO
	//defer triggerGenNewsfeed(post.UserId)

	c.JSON(http.StatusOK, gin.H{
		"postId": postId,
	})
}

func (app *Handler) UploadImage(c *gin.Context) {
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
	imagePath, err := app.imageStorage.Put(out, postId, filename, fStat.Size())
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	//err = app.postDB.(postId, imagePath)
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

func (app *Handler) EditPost(c *gin.Context) {
	var newPost payloads.PostPayload
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

	curPost, err := app.postDB.GetPostById(postId)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	curPost.Merge(&newPost)
	err = app.postDB.Edit(postId, &curPost)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	// TODO
	//defer triggerGenNewsfeed(curPost.UserId)

	// TODO Trigger delete cache if possible
	c.JSON(http.StatusOK, gin.H{
		"message": "Need implementation",
	})
}

func (app *Handler) DeletePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = app.postDB.Delete(postId)
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

func (app *Handler) LikePost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var p payloads.PostPayload
	err = c.ShouldBindBodyWithJSON(&p)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	err = app.postDB.Like(postId, p.UserId)
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

func (app *Handler) CommentPost(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("post_id"))

	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	var cmt payloads.CommentRecord
	err = c.ShouldBindBodyWithJSON(&cmt)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "err: " + err.Error(),
		})
		return
	}

	cmtId, err := app.postDB.Comment(postId, cmt.UserId, cmt.ContentText)
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

// RouteV1 Route for post service
func RouteV1(h *Handler, r *gin.Engine) {
	v1 := r.Group("v1")
	v1.Use()
	{
		v1.GET("/posts/:post_id", h.GetPost)
		v1.POST("/posts", h.CreatePost)
		v1.POST("/posts/:post_id/images", h.UploadImage)
		v1.PUT("/posts/:post_id", h.EditPost)
		v1.DELETE("/posts/:post_id", h.DeletePost)
		v1.POST("/posts/:post_id/comments", h.CommentPost)
		v1.POST("/posts/:post_id/likes", h.LikePost)
	}
}

func NewHandler(postDB database.PostDB, imageStorage object_store.ImageStorage) *Handler {
	return &Handler{
		postDB:       postDB,
		imageStorage: imageStorage,
	}
}
