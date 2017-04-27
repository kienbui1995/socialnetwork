package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/models"
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
		status := models.UserStatus{UserID: userid, Message: json.Message}
		status.ID, erruid = services.CreateUserStatus(status)
		if erruid == nil && status.ID >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Create user status successful", map[string]interface{}{"id": status.ID})
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
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		statusList, errList := services.GetUserStatuses(userid)
		if errList == nil && statusList != nil {
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

		status, errUpdate := services.UpdateUserStatus(statusid, json.Message)
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

// // UpdatePost func to update info a Post
// func UpdatePost(c *gin.Context) {
// 	var (
// 		err     error
// 		postid  int
// 		content string
// 		image   string
// 		status  int
// 		update  bool
// 	)
// 	postid, err = strconv.Atoi(c.Param("postid"))
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"code":    -1,
// 			"message": err.Error(),
// 			"postid":  postid,
// 		})
// 		return
// 	}
// 	content = c.PostForm("content")
// 	image = c.PostForm("image")
// 	status, err = strconv.Atoi(c.DefaultPostForm("status", ""))
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"code":    -1,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	post := models.Post{}
// 	post.PostID = postid
// 	post.Content = content
// 	post.Image = image
// 	post.Status = status
// 	post.UpdatedTime = time.Now().String()
// 	fmt.Printf("%s wsssss", post.UpdatedTime)
// 	update, err = services.UpdatePost(post)
// 	if err != nil {
// 		c.JSON(200, gin.H{
// 			"code":    -1,
// 			"message": err.Error(),
// 		})
// 		return
// 	}
// 	if update == true {
// 		c.JSON(200, gin.H{
// 			"code":    1,
// 			"message": "Uppdate post successful",
// 			"postid":  post.PostID,
// 		})
// 	} else {
// 		c.JSON(200, gin.H{
// 			"code":    -1,
// 			"message": "Don't update info in DB",
// 		})
// 	}
//
// }
//
// // DeletePost func to delete a Post
// func DeletePost(c *gin.Context) {
// }
