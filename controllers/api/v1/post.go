package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/models"
)

// GetPost func  return info Post
func GetPost(c *gin.Context) {
	c.JSON(200, gin.H{
		"Sss": "Sss",
	})
}

// CreatePost func to create a new post
func CreatePost(c *gin.Context) {
	post := models.Post{}
	post.Message = c.DefaultPostForm("content", "")
	post.IsHidden = false
	post.Status = 1
	post.Type = c.DefaultPostForm("type", "post")
}

// UpdatePost func to update info a Post
func UpdatePost(c *gin.Context) {
}

// DeletePost func to delete a Post
func DeletePost(c *gin.Context) {
}
