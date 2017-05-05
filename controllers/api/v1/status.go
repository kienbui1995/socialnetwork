package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/services"
)

// CreateUserStatus func  to create a user status
func CreateUserStatus(c *gin.Context) {
	userid, erruid := strconv.Atoi(c.Param("userid"))
	if erruid != nil {
		libs.ResponseBadRequestJSON(c, 110, "Invalid user id")
	} else {

		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		json := struct {
			Message string `json:"message"`
			Privacy string `json:"privacy"`
			Status  string `json:"status"`
		}{}
		if errBind := c.Bind(&json); errBind != nil {
			libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
			return
		}

		// validation
		if len(json.Message) == 0 {
			libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
			return
		}

		statusID, errsid := services.CreateUserStatus(userid, json.Message, json.Privacy, json.Status)
		if errsid == nil && statusID >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Create user status successful", map[string]interface{}{"id": statusID})

			// auto Increase Posts
			go func() {
				ok, err := services.IncreasePosts(userid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePosts service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePosts service")
				}
			}()

			// push noti
			go func() {
				ids, errGetIDs := services.GetFollowerIDs(userid)
				if len(ids) > 0 && errGetIDs == nil {
					for index := 0; index < len(ids); index++ {
						PushTest(ids[index], 1, "Một người bạn follow viết bài", json.Message)
					}

				}
			}()

			return
		}
		libs.ResponseServerErrorJSON(c)
		if erruid != nil {
			fmt.Printf("ERROR in CreateUserStatus services: %s", erruid.Error())
		} else {
			fmt.Printf("ERROR in CreateUserStatus services: Don't create User Status")
		}
	}
}

// GetUserStatuses func to create a new post
func GetUserStatuses(c *gin.Context) {
	userid, erruid := strconv.Atoi(c.Param("userid"))
	if erruid != nil {
		libs.ResponseBadRequestJSON(c, 110, "Invalid user id")
	} else {

		//check permisson
		myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		sort := c.DefaultQuery("sort", "-created_at")
		print(sort)
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

		statusList, errList := services.GetUserStatuses(userid, myuserid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Statuses List", statusList, nil, len(statusList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetUserStatuses services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetUserStatuses services: Don't get User Statuses")
		}

	}
}

// UpdateUserStatus func to create a new post
func UpdateUserStatus(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid status id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserStatus(statusid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDPostedStatus(statusid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		json := struct {
			Message string `json:"message" valid:"nonzero"`
			Privacy int    `json:"privacy"`
			Status  int    `json:"status"`
		}{}
		if errBind := c.Bind(&json); errBind != nil {
			libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
			return
		}

		// validation
		if len(json.Message) == 0 {
			libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
			return
		}

		status, errUpdate := services.UpdateUserStatus(statusid, json.Message, json.Privacy, json.Status)
		if errUpdate == nil && status.CreatedAt > 0 {
			libs.ResponseSuccessJSON(c, 1, "Update user status successful", status)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errUpdate != nil {
			fmt.Printf("ERROR in UpdateUserStatus services: %s", errUpdate.Error())
		} else {
			fmt.Printf("ERROR in UpdateUserStatus services: Don't update User Statuses")
		}

	}
}

// DeleteUserStatus func to delete a user status
func DeleteUserStatus(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid status id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserStatus(statusid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDPostedStatus(statusid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		ok, errDel := services.DeleteUserStatus(statusid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Delete user status successful", nil)

			// auto Decrease Posts
			go func() {
				ok, err := services.DecreasePosts(userid)
				if err != nil {
					fmt.Printf("ERROR in DecreasePosts service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in DecreasePosts service")
				}
			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errDel != nil {
			fmt.Printf("ERROR in DeleteUserStatus services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeleteUserStatus services: Don't delete User Statuses")
		}

	}
}

// GetUserStatus func to delete a user status
func GetUserStatus(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid status id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserStatus(statusid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		// userid, _ := services.GetUserIDPostedStatus(statusid)
		//check permisson ~needfix when privacy not public
		myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if myuserid == -1 || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		status, errGet := services.GetUserStatus(statusid, myuserid)
		if errGet == nil && status.ID == statusid {
			libs.ResponseSuccessJSON(c, 1, "Get user status successful", status)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errGet != nil {
			fmt.Printf("ERROR in GetUserStatus services: %s", errGet.Error())
		} else {
			fmt.Printf("ERROR in GetUserStatus services: Don't get User Status")
		}

	}
}

// CreateStatusLike func
func CreateStatusLike(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// check liked
		if liked, _ := services.CheckExistStatusLike(statusid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		liked, errLike := services.CreateStatusLike(statusid, userid)
		if errLike == nil && liked == true {
			libs.ResponseSuccessJSON(c, 1, "Like status successful", nil)

			// auto Increase Status Likes
			go func() {
				ok, err := services.IncreaseStatusLikes(statusid)
				if err != nil {
					fmt.Printf("ERROR in IncreaseStatusLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreaseStatusLikes service")
				}
			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreateStatusLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreateStatusLike services: Don't CreateStatusLike")
		}

	}
}

// DeleteStatusLike func to delete a like
func DeleteStatusLike(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserStatus(statusid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		//check permisson
		userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		// check liked
		if liked, _ := services.CheckExistStatusLike(statusid, userid); liked != true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		ok, errDel := services.DeleteStatusLike(statusid, userid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Unlike successful", nil)

			// auto Decrease Status Likes
			go func() {
				ok, err := services.DecreaseStatusLikes(statusid)
				if err != nil {
					fmt.Printf("ERROR in DecreaseStatusLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in DecreaseStatusLikes service")
				}
			}()

			return
		}

		libs.ResponseServerErrorJSON(c)
		if errDel != nil {
			fmt.Printf("ERROR in DeleteStatusLike services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeleteStatusLike services: Don't DeleteStatusLike")
		}

	}
}

// GetStatusLikes func
func GetStatusLikes(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		sort := c.DefaultQuery("sort", "-liked_at")

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

		likeList, errList := services.GetStatusLikes(statusid, userid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, " Statuses Likes User List", likeList, nil, len(likeList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetStatusLikes services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetStatusLikes services: Don't GetStatusLikes")
		}

	}
}
