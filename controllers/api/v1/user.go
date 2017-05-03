package v1

import (
	"fmt"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetUser func  return info user
func GetUser(c *gin.Context) {
	if len(c.Param("userid")) == 0 {
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

	if govalidator.IsByteLength(user.Username, 3, 15) == false {
		libs.ResponseErrorJSON(c, libs.NewErrorDetail(382, "Please enter a valid username."))
		return
	}
	if govalidator.IsEmail(user.Email) == false {
		libs.ResponseErrorJSON(c, libs.NewErrorDetail(385, "Please enter a valid email address."))
		return
	}

	if exist, _ := services.CheckExistUsername(user.Username); exist == true {
		libs.ResponseErrorJSON(c, libs.NewErrorDetail(376, "The login credential you provided belongs to an existing account"))
		return
	}

	if exist, _ := services.CheckExistEmail(user.Email); exist == true {
		libs.ResponseErrorJSON(c, libs.NewErrorDetail(371, "The email address you provided belongs to an existing account"))
		return
	}

	user.Status = 0
	activecode := libs.RandStringBytes(6)
	go func() {
		if err := services.CreateEmailActive(models.EmailActive{Email: user.Email, ActiveCode: activecode, UserID: user.UserID}); err != nil {
			fmt.Printf("\nCreate Email Active faile: %s", err.Error())
		}

	}()
	// /fmt.Printf("%v", user.EmailActive)
	user.UserID, errUser = services.CreateUserTest(user) //needfix
	if errUser != nil {
		libs.ResponseJSON(c, 400, 387, "There was an error with your registration. Please try registering again: "+errUser.Error(), nil)
		return

	}

	// pause func send mail active
	// go func() {
	// 	sender := libs.NewSender("kien.laohac@gmail.com", "ytanyybkizzygqjk")
	// 	var email []string
	// 	email = append(email, user.Email)
	// 	linkActive := "<a href='localhost:8080/user/" + string(user.UserID) + "?email_active=" + activecode + "'>Active</a>"
	// 	sender.SendMail(email, fmt.Sprintf("Active user %s on TLSEN", user.Username), fmt.Sprintf("Content-Type: text/html; charset=UTF-8\n\ncode: %s OR active via link: %s", activecode, linkActive))
	// }()

	libs.ResponseCreatedJSON(c, 1, "Create user successful!", map[string]interface{}{"id": user.UserID})

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
	libs.ResponseCreatedJSON(c, 1, "Create user successful!", map[string]interface{}{"id": user.UserID})

}

// UpdateUser func to update info a User
func UpdateUser(c *gin.Context) {
	var update bool

	userid, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		libs.ResponseBadRequestJSON(c, 110, "Invalid user id: "+err.Error())
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": err.Error(),
		// })
		return
	}
	var jsonUser map[string]interface{}

	if c.Bind(&jsonUser) != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter"+c.Bind(&jsonUser).Error())
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": err.Error(),
		// })
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
		libs.ResponseBadRequestJSON(c, 310, "User data edit failure")
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": err.Error(),
		// })
		return
	}
	if update == true {
		c.JSON(200, gin.H{
			"code":    1,
			"message": "Update successful",
			"data":    user,
		})
	} else {
		libs.ResponseBadRequestJSON(c, 310, "User data edit failure")
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": "Don't update info in DB",
		// })
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

//CreateUserSubscribers func
func CreateUserSubscribers(c *gin.Context) {
	var sub = models.Subscriber{}
	var errSub error
	var err error
	//var json interface{}

	if len(c.Param("userid")) == 0 || len(c.Param("toid")) == 0 {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: userid", nil)
		return
	}

	if sub.FromID, err = strconv.Atoi(c.Param("userid")); err != nil {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: userid", nil)
		return
	}
	if sub.ToID, err = strconv.Atoi(c.Param("toid")); err != nil {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: toid", nil)
		return
	}
	if exist1, _ := services.CheckExistNodeWithID(sub.FromID); exist1 != true {
		libs.ResponseJSON(c, 400, 110, "Invalid user id", nil)
		return
	}

	//check permisson
	if userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != sub.FromID || errGet != nil {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	if exist2, _ := services.CheckExistNodeWithID(sub.ToID); exist2 != true {
		libs.ResponseJSON(c, 400, 110, "Invalid user id", nil)
		return
	}

	sub.SubscriberID, errSub = services.CreateUserSubscriber(sub.FromID, sub.ToID)
	if errSub != nil {
		libs.ResponseJSON(c, 400, -1, errSub.Error(), nil)
		return

	}

	libs.ResponseCreatedJSON(c, 1, "Create subscriber successful!", sub)
	device, _ := services.GetDeviceByUserID(sub.ToID)
	PushNotification(device)
	// auto Increase Followers And Followings
	go func() {
		ok, err := services.IncreaseFollowersAndFollowings(sub.FromID, sub.ToID)
		if err != nil {
			fmt.Printf("ERROR in IncreaseFollowersAndFollowings service: %s", err.Error())
		}
		if ok != true {
			fmt.Printf("ERROR in IncreaseFollowersAndFollowings service")
		}
	}()
}

//GetFollowers func
func GetFollowers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: userid")
		return
	}
	check, errCheck := CheckPermissionsWithID(id, c.Request.Header.Get("token"))
	if errCheck != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR in CheckPermissionsWithID: %s", errCheck.Error())
		return
	}

	if check == false {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	SUserList, errGet := services.GetFollowers(id)
	if errGet != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR in GetFollowers: %s", errGet.Error())
		return
	}
	libs.ResponseEntityListJSON(c, 1, "User list", SUserList, nil, len(SUserList))
}

//GetSubscribers func
func GetSubscribers(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("userid"))
	if err != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: userid")
		return
	}
	check, errCheck := CheckPermissionsWithID(id, c.Request.Header.Get("token"))
	if errCheck != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR in CheckPermissionsWithID: %s", errCheck.Error())
		return
	}

	if check == false {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}
	SUserList, errGet := services.GetSubscribers(id)
	if errGet != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR in GetSubscribers: %s", errGet.Error())
		return
	}
	libs.ResponseEntityListJSON(c, 1, "User list", SUserList, nil, len(SUserList))
}

//DeleteUserSubscribers func
func DeleteUserSubscribers(c *gin.Context) {
	var sub = models.Subscriber{}

	//var json interface{}

	if len(c.Param("userid")) == 0 || len(c.Param("toid")) == 0 {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: userid", nil)
		return
	}

	sub.FromID, _ = strconv.Atoi(c.Param("userid"))
	sub.ToID, _ = strconv.Atoi(c.Param("toid"))
	if exist1, _ := services.CheckExistNodeWithID(sub.FromID); exist1 != true {
		libs.ResponseJSON(c, 400, 110, "Invalid user id", nil)
		return
	}

	//check permission
	if userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != sub.FromID || errGet != nil {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
	}
	if exist2, _ := services.CheckExistNodeWithID(sub.ToID); exist2 != true {
		libs.ResponseJSON(c, 400, 110, "Invalid user id", nil)
		return
	}

	if exist, err := services.CheckExistUserSubscriber(sub.FromID, sub.ToID); exist == false || err != nil {
		libs.ResponseJSON(c, 400, 2, "No exist this object", nil)
		return
	}

	delsub, errdel := services.DeleteUserSubscriber(sub.FromID, sub.ToID)
	if errdel == nil && delsub == true {

		// auto Decrease Followers And Followings
		go func() {
			ok, err := services.DecreaseFollowersAndFollowings(sub.FromID, sub.ToID)
			if err != nil {
				fmt.Printf("ERROR in DecreaseFollowersAndFollowings service: %s", err.Error())
			}
			if ok != true {
				fmt.Printf("ERROR in DecreaseFollowersAndFollowings service")
			}
		}()

		libs.ResponseSuccessJSON(c, 1, "Delete subscriber successful", nil)
		return

	}
	if errdel != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR in DeleteUserSubscriber: %s", errdel.Error())
	}

}

//FindUser func
func FindUser(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		libs.ResponseBadRequestJSON(c, 101, "Missing a few fields: name")
		return
	}
	token := c.Request.Header.Get("token")
	id, _ := GetUserIDFromToken(token)
	userList, errFind := services.FindUserByUsernameAndFullName(id, name)
	if errFind != nil {
		libs.ResponseServerErrorJSON(c)
		fmt.Printf("ERROR FindUserByUsernameAndFullName: %s", errFind.Error())
	} else {
		libs.ResponseEntityListJSON(c, 1, "User list found", userList, nil, len(userList))
	}

}

// CreateUserTest func to create a new user
func CreateUserTest(c *gin.Context) {
	var user = models.User{}
	var errUser error
	//var json interface{}

	if c.Bind(&user) != nil {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+c.Bind(&user).Error(), nil)
		return
	}

	user.UserID, errUser = services.CreateUserTest(user)
	if errUser != nil {
		libs.ResponseJSON(c, 400, -1, errUser.Error(), nil)
		return

	}
	libs.ResponseCreatedJSON(c, 1, "Create user successful!", map[string]interface{}{"id": user.UserID})

}

// GetNewsFeed func to create a new post
func GetNewsFeed(c *gin.Context) {
	userid, erruid := strconv.Atoi(c.Param("userid"))
	if erruid != nil {
		libs.ResponseBadRequestJSON(c, 110, "Invalid user id")
	} else {

		//check permisson
		myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if userid != myuserid || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		sort := c.DefaultQuery("sort", "-created_at")
		orderby, errSort := libs.ConvertSort(sort)
		if errSort != nil {
			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSort.Error())
			return
		}
		skip, errSkip := strconv.Atoi(c.DefaultQuery("skip", "0"))
		if errSkip != nil {
			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSkip.Error())
			return
		}
		limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "25"))
		if errLimit != nil {
			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errLimit.Error())
			return
		}

		statusList, errList := services.GetNewsFeed(userid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "Get news feed successful", statusList, nil, len(statusList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetNewsFeed services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetNewsFeed services: Don't get GetNewsFeed")
		}

	}
}
