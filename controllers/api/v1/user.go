package v1

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
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
		type Paging struct {
			NextStart     int `json:"next_start"`
			PreviousStart int `json:"previous_start"`
		}
		var paging Paging
		var next, previous int
		if limit != 0 || start != 0 {
			if limit == 0 {
				limit = 25
			}
			listuser, errlist = services.GetAllUserWithSkipLimit(start, limit)
			next = start + len(listuser)
			if start-len(listuser) >= 0 {
				previous = start - len(listuser)
			} else {
				previous = 0
			}
		} else {
			listuser, errlist = services.GetAllUser()
			next = len(listuser)
			previous = 0
		}
		paging = Paging{next, previous}
		if errlist != nil {
			libs.ResponseJSON(c, 200, -1, errlist.Error(), nil)
		} else {
			libs.ResponseEntityListJSON(c, 1, "list user", listuser, paging, len(listuser))
			return
		}
	} else {
		userid, err := strconv.Atoi(c.Param("userid"))
		if err != nil {
			libs.ResponseBadRequestJSON(c, 110, "Invalid user id"+err.Error()) // NeedEdit
		} else {
			var errUser error
			var user models.User

			// user.UserID = userid
			user, errUser = services.GetUser(userid)
			if errUser != nil {
				libs.ResponseNotFoundJSON(c, -1, errUser.Error())
			} else {
				libs.ResponseSuccessJSON(c, 1, "info a user", user)
			}
		}
	}
}

// SignUp func stead of CreateUser func but  active  = 0 and has email Active
func SignUp(c *gin.Context) {
	var user = models.User{}
	var errUser error
	//var json interface{}

	errBind := c.BindJSON(&user)
	if errBind != nil {
		libs.ResponseBadRequestJSON(c, -1, errBind.Error())
		return
	}

	errorDetails := []libs.ErrorDetail{}
	if govalidator.IsByteLength(user.Username, 3, 15) == false {
		errorDetails = append(errorDetails, libs.NewErrorDetail(382, "Please enter a valid username."))
	}
	if govalidator.IsEmail(user.Email) == false {
		errorDetails = append(errorDetails, libs.NewErrorDetail(385, "Please enter a valid email address."))
	}
	if len(errorDetails) != 0 {
		libs.ResponseErrorJSON(c, libs.Errors{Code: 387, Message: "There was an error with your registration. Please try registering again."})
		return
	}
	user.Status = 0
	user.EmailActive = libs.RandStringBytes(6)
	fmt.Printf("%v", user.EmailActive)
	user.UserID, errUser = services.CreateUser(user)
	if errUser != nil {
		libs.ResponseJSON(c, 400, -1, errUser.Error(), nil)
		return

	}
	sender := libs.NewSender("kien.laohac@gmail.com", "ytanyybkizzygqjk")
	var email []string
	email = append(email, user.Email)
	linkActive := "<a href='localhost:8080/user/" + string(user.UserID) + "?email_active=" + user.EmailActive + "'>Active</a>"
	go sender.SendMail(email, fmt.Sprintf("Active user %s on TLSEN", user.Username), fmt.Sprintf("Content-Type: text/html; charset=UTF-8\n\ncode: %s OR active via link: %s", user.EmailActive, linkActive))

	libs.ResponseCreatedJSON(c, 1, "Create user successful!", user.UserID)

}

// CreateUser func to create a new user
func CreateUser(c *gin.Context) {
	var user = models.User{}
	var errUser error
	//var json interface{}

	if c.Bind(&user) != nil {
		libs.ResponseJSON(c, 400, -1, c.Bind(&user).Error(), nil)
		return
	}

	user.UserID, errUser = services.CreateUser(user)
	if errUser != nil {
		libs.ResponseJSON(c, 400, -1, errUser.Error(), nil)
		return

	}
	libs.ResponseCreatedJSON(c, 1, "Create user successful!", user.UserID)

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
