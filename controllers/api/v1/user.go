package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetUser func  return info user
func GetUser(c *gin.Context) {
	//fmt.Printf("query: %s", c.Param("userid"))
	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
	} else {
		var errUser error
		user := models.User{}
		user.UserID = userid
		user, errUser = services.GetUser(user.UserID)
		if errUser != nil {
			c.JSON(200, gin.H{
				"error": errUser.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"user": user,
			})
		}
	}

}

// CreateUser func to create a new user
func CreateUser(c *gin.Context) {
	var user = models.User{}
	var errUser error
	userid, _ := strconv.Atoi(c.DefaultPostForm("userid", "-1"))
	username := c.DefaultPostForm("username", "")
	password := c.DefaultPostForm("password", "")
	email := c.DefaultPostForm("email", "")
	if username != "" && password != "" && email != "" {
		user.UserID = int(userid)
		user.Username = username
		user.Password = password
		user.Email = email
		user.Status = 0
	} else {
		c.JSON(200, gin.H{
			"idUser": -1,
			"error":  "Missing a few post value",
		})
		return
	}
	user.UserID, errUser = services.CreateUser(user)
	if errUser != nil {
		c.JSON(200, gin.H{
			"userid":  user.UserID,
			"error":   errUser,
			"message": "Created new user",
		})
		return

	}
	c.JSON(200, gin.H{
		"userid": user.UserID,
	})

}

// UpdateUser func to update info a User
func UpdateUser(c *gin.Context) {
	var (
		err    error
		userid int
		email  string
		status int
		update bool
	)
	userid, err = strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	email = c.PostForm("email")
	status, err = strconv.Atoi(c.PostForm("status"))
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{}
	user.UserID = userid
	user.Email = email
	user.Status = status
	update, err = services.UpdateUser(user)
	if err != nil {
		c.JSON(200, gin.H{
			"error": err.Error(),
		})
		return
	}
	if update == true {
		c.JSON(200, gin.H{
			"userid": user.UserID,
		})
	} else {
		c.JSON(200, gin.H{
			"userid": -1,
		})
	}

}

// DeleteUser func to delete a User
func DeleteUser(c *gin.Context) {
	userid, err1 := strconv.Atoi(c.Param("userid"))
	if err1 != nil {
		c.JSON(200, gin.H{
			"err1": err1,
		})
		return
	}
	del, err := services.DeleteUser(userid)

	c.JSON(200, gin.H{
		"deleted": del,
		"err":     err,
	})

}
