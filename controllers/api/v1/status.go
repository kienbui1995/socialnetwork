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
			Message string `json:"message" valid:"nonzero"`
			Privacy int    `json:"privacy"`
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

		statusID, errsid := services.CreateUserStatus(userid, json.Message, json.Privacy, 1)
		if errsid == nil && statusID >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Create user status successful", map[string]interface{}{"id": statusID})
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

		// //check permisson
		// if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
		// 	libs.ResponseAuthJSON(c, 200, "Permissions error")
		// 	return
		// }

		sort := c.DefaultQuery("sort", "+created_at")

		orderby := libs.ConvertSort(sort)
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

		statusList, errList := services.GetUserStatuses(userid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Statuses List", statusList, nil, len(statusList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if erruid != nil {
			fmt.Printf("ERROR in GetUserStatuses services: %s", erruid.Error())
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
		// if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
		// 	libs.ResponseAuthJSON(c, 200, "Permissions error")
		// 	return
		// }

		status, errGet := services.GetUserStatus(statusid)
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

		// // binding
		// json := struct {
		// 	Message string `json:"message" valid:"nonzero"`
		// }{}
		// if errBind := c.Bind(&json); errBind != nil {
		// 	libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
		// 	return
		// }

		// // validation
		// if len(json.Message) == 0 {
		// 	libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
		// 	return
		// }

		liked, errLike := services.CreateStatusLike(statusid, userid)
		if errLike == nil && liked == true {
			libs.ResponseSuccessJSON(c, 1, "Like status successful", nil)
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

		ok, errDel := services.DeleteStatusLike(statusid, userid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Unlike successful", nil)
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

		sort := c.DefaultQuery("sort", "+created_at")

		orderby := libs.ConvertSort(sort)
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
