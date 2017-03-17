package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetPost func  return info Post
func GetPost(c *gin.Context) {
	userid, erruid := strconv.Atoi(c.Param("userid"))
	if erruid != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": erruid.Error(),
		})
	} else {
		listpost, errlist := services.GetPostByUserID(userid)
		if errlist != nil {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": errlist.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "List post by userid",
				"data":    listpost,
			})
		}
	}

}

// CreatePost func to create a new post
func CreatePost(c *gin.Context) {
	post := models.Post{}
	post.Content = c.DefaultPostForm("content", "")
	post.Image = c.DefaultPostForm("image", "")
	post.Status = 1
	post.CreatedTime = c.DefaultPostForm("createdtime", time.Now().String())
	post.UpdatedTime = c.DefaultPostForm("updatedtime", time.Now().String())
	userid, errid := strconv.Atoi(c.Param("userid"))
	if errid != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": errid.Error(),
		})
	}
	post, err := services.CreatePost(post, userid)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "info a post",
			"data":    post,
		})
	}
}

// UpdatePost func to update info a Post
func UpdatePost(c *gin.Context) {
	var (
		err     error
		postid  int
		content string
		image   string
		status  int
		update  bool
	)
	postid, err = strconv.Atoi(c.Param("postid"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
			"postid":  postid,
		})
		return
	}
	content = c.PostForm("content")
	image = c.PostForm("image")
	status, err = strconv.Atoi(c.DefaultPostForm("status", ""))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	post := models.Post{}
	post.PostID = postid
	post.Content = content
	post.Image = image
	post.Status = status
	post.UpdatedTime = time.Now().String()
	fmt.Printf("%s wsssss", post.UpdatedTime)
	update, err = services.UpdatePost(post)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	if update == true {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "Uppdate post successful",
			"postid":  post.PostID,
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Don't update info in DB",
		})
	}

}

// DeletePost func to delete a Post
func DeletePost(c *gin.Context) {
}
