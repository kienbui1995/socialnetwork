package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetUser func  return info user
func GetUser(c *gin.Context) {
	if c.Param("userid") == "" {
		limit, _ := strconv.Atoi(c.Query("limit"))
		start, _ := strconv.Atoi(c.Query("start"))
		var listuser []models.User
		var errlist error
		if limit != 0 && start != 0 {
			listuser, errlist = services.GetAllUserWithSkipLimit(start, limit)
		} else {
			listuser, errlist = services.GetAllUser()
		}
		if errlist != nil {
			libs.ResponseJSON(c, 200, -1, errlist.Error(), nil)
			// c.JSON(200, gin.H{
			// 	"code":    -1,
			// 	"message": errlist.Error(),
			// })
		} else {
			libs.ResponseEntityListJSON(c, 1, "list user", listuser, nil, len(listuser))
			return
		}
	} else {
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
}

// CreateUser func to create a new user
func CreateUser(c *gin.Context) {
	var user = models.User{}
	var errUser error
	//var json interface{}

	if c.Bind(&user) != nil {
		libs.ResponseJSON(c, 400, -1, c.Bind(&user).Error(), nil)
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": c.Bind(&user).Error(),
		// })
		return
	}
	// errCopy := copier.Copy(&user, &json)
	//fmt.Printf("\n%v", user)
	// // fmt.Printf("\n%v", json)
	// if errCopy != nil {
	// 	libs.ResponseJSON(c, 200, -1, errCopy.Error(), nil)
	// 	return
	// }

	user.UserID, errUser = services.CreateUser(user)
	if errUser != nil {
		libs.ResponseJSON(c, 400, -1, errUser.Error(), nil)
		// c.JSON(200, gin.H{
		// 	"code": -1,
		// 	//"userid":  user.UserID,
		// 	"message": errUser.Error(),
		// 	//"message": "Created new user",
		// })
		return

	}
	libs.ResponseCreatedJSON(c, 1, "Create user successful!", user.UserID)
	// c.JSON(200, gin.H{
	// 	"code":    1,
	// 	"message": "Create user successful!",
	// 	"data":    user.UserID,
	// })

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
	user.UserID = userid
	errCopy := copier.Copy(user, jsonUser)
	if errCopy != nil {
		libs.ResponseJSON(c, 200, 1, "Thong bao moi.", nil)
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
			"userid":  user.UserID,
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
