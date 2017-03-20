package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetUser func  return info user
func GetUser(c *gin.Context) {
	if c.Param("userid") == "" {
		listuser, errlist := services.GetAllUser()
		if errlist != nil {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": errlist.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "list user",
				"data":    listuser,
			})
			return
		}
	}
	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
	} else {
		var errUser error
		var user models.User
		// user.UserID = userid
		user, errUser = services.GetUser(userid)
		if errUser != nil {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": errUser.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "info a user",
				"data":    user,
			})
		}
	}

}

// CreateUser func to create a new user
func CreateUser(c *gin.Context) {
	var user = models.User{}
	var errUser error
	var json map[string]interface{}

	if c.Bind(&json) != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": c.Bind(&json).Error(),
		})
		return
	}
	user.Data = make(map[string]interface{})
	for k, v := range json {
		user.Data[k] = v
	}

	user, errUser = services.CreateUser(user)
	if errUser != nil {
		c.JSON(200, gin.H{
			"code": -1,
			//"userid":  user.UserID,
			"message": errUser.Error(),
			//"message": "Created new user",
		})
		return

	}

	c.JSON(200, gin.H{
		"code":    1,
		"message": "Create user successful!",
		"data":    user.Data,
	})

}

// UpdateUser func to update info a User
func UpdateUser(c *gin.Context) {
	var update bool

	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	var jsonUser map[string]interface{}

	if c.Bind(&jsonUser) != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	var user models.User
	user.Data = make(map[string]interface{})
	user.Data["userid"] = userid
	for k, v := range jsonUser {
		user.Data[k] = v
	}
	update, err = services.UpdateUser(user)
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
			"message": "Update successful",
			"userid":  user.Data["userid"],
		})
	} else {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Don't update info in DB",
		})
	}

}

// DeleteUser func to delete a User
func DeleteUser(c *gin.Context) {
	userid, err1 := strconv.Atoi(c.Param("userid"))
	if err1 != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Don't convert param userid",
		})
		return
	}
	del, err := services.DeleteUser(userid)

	c.JSON(200, gin.H{
		"code":    1,
		"deleted": del,
		"err":     err,
	})

}
