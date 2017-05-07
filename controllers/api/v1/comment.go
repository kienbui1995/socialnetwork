package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetComments func
func GetComments(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
		if allowed, _ := services.CheckPostInteractivePermission(postid, userid); allowed == false {
			libs.ResponseForbiddenJSON(c, 220, "Post not visible")
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

		commentList, errList := services.GetComments(postid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Post Comments List", commentList, nil, len(commentList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetComments services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetComments services: Don't GetComments")
		}

	}
}

// GetComment func to get a comment
func GetComment(c *gin.Context) {
	commentid, errcid := strconv.Atoi(c.Param("commentid"))
	if errcid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistComment(commentid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		// userid, _ := services.GetUserIDPostedStatus(postid)
		//check permisson ~needfix when privacy not public
		// myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		// if myuserid == -1 || errGet != nil {
		// 	libs.ResponseAuthJSON(c, 200, "Permissions error")
		// 	return
		// }

		comment, errGet := services.GetComment(commentid)
		if errGet == nil && comment.ID == commentid {
			libs.ResponseSuccessJSON(c, 1, "Get comment successful", comment)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errGet != nil {
			fmt.Printf("ERROR in GetComment services: %s", errGet.Error())
		} else {
			fmt.Printf("ERROR in GetComment services: Don't GetComment")
		}

	}
}

// CreateComment func  to create comment
func CreateComment(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	}
	//check permisson
	userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
	if errGet != nil {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	json := struct {
		Message string `json:"message"`
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
	// validation

	if json.Status == 0 {
		json.Status = 1
	}

	commentID, errcid := services.CreateComment(userid, json.Message, json.Status, postid)
	if errcid == nil && commentID >= 0 {
		libs.ResponseSuccessJSON(c, 1, "Create comment successful", map[string]interface{}{"id": commentID})

		// auto Increase Posts
		go func() {
			ok, err := services.IncreaseObjectComments(postid)
			if err != nil {
				fmt.Printf("ERROR in IncreaseObjectComments service: %s", err.Error())
			}
			if ok != true {
				fmt.Printf("ERROR in IncreaseObjectComments service")
			}
		}()

		// push noti
		go func() {
			user, _ := services.GetUser(userid)
			id, _ := services.GetUserIDByPostID(postid)

			if id >= 0 {
				PushTest(id, 1, "@"+user.Username+" bình luận bài đăng của bạn", json.Message)
			}

		}()
		return
	}
	libs.ResponseServerErrorJSON(c)
	if errcid != nil {
		fmt.Printf("ERROR in CreateComment services: %s", errcid.Error())
	} else {
		fmt.Printf("ERROR in CreateComment services: Don't CreateComment")
	}

}

// UpdateComment func
func UpdateComment(c *gin.Context) {
	commentid, errcid := strconv.Atoi(c.Param("commentid"))
	if errcid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistComment(commentid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDWroteComment(commentid)
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

		updated, errUpdate := services.UpdateComment(commentid, json.Message)
		if errUpdate == nil && updated == true {
			libs.ResponseSuccessJSON(c, 1, "Update comment successful", nil)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errUpdate != nil {
			fmt.Printf("ERROR in UpdateComment services: %s", errUpdate.Error())
		} else {
			fmt.Printf("ERROR in UpdateComment services: Don't UpdateComment")
		}
	}
}

// DeleteComment func to delete a comment
func DeleteComment(c *gin.Context) {
	commentid, errcid := strconv.Atoi(c.Param("commentid"))
	if errcid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistComment(commentid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDWroteComment(commentid)
		objectid, _ := services.GetObjectIDbyCommentID(commentid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		ok, errDel := services.DeleteComment(commentid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Delete comment successful", nil)

			// auto Decrease Status Comments
			go func() {
				ok, err := services.DecreaseObjectComments(objectid)
				if err != nil {
					fmt.Printf("ERROR in IncreaseStatusComments service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreaseStatusComments service")
				}
			}()

			return
		}

		libs.ResponseServerErrorJSON(c)
		if errDel != nil {
			fmt.Printf("ERROR in DeleteStatusComment services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeleteStatusComment services: Don't DeleteStatusComment")
		}

	}
}

// // GetStatusComments func
// func GetStatusComments(c *gin.Context) {
// 	statusid, errsid := strconv.Atoi(c.Param("statusid"))
// 	if errsid != nil {
// 		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
// 	} else {
//
// 		//check permisson
// 		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
// 		if allowed, _ := services.CheckStatusInteractivePermission(statusid, userid); allowed == false {
// 			libs.ResponseForbiddenJSON(c, 220, "Status not visible")
// 			return
// 		}
//
// 		sort := c.DefaultQuery("sort", "-created_at")
// 		print(sort)
// 		orderby, errSort := libs.ConvertSort(sort)
// 		if errSort != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSort.Error())
// 			return
// 		}
// 		skip, errSkip := strconv.Atoi(c.DefaultQuery("skip", "0"))
// 		if errSkip != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSkip.Error())
// 			return
// 		}
// 		limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "25"))
// 		if errLimit != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errLimit.Error())
// 			return
// 		}
//
// 		commentList, errList := services.GetStatusComments(statusid, orderby, skip, limit)
// 		if errList == nil {
// 			libs.ResponseEntityListJSON(c, 1, "User Statuses Comments List", commentList, nil, len(commentList))
// 			return
// 		}
//
// 		libs.ResponseServerErrorJSON(c)
// 		if errList != nil {
// 			fmt.Printf("ERROR in GetStatusComments services: %s", errList.Error())
// 		} else {
// 			fmt.Printf("ERROR in GetStatusComments services: Don't getGetStatusComments")
// 		}
//
// 	}
// }
//
// // CreateStatusComment func
// func CreateStatusComment(c *gin.Context) {
// 	statusid, errsid := strconv.Atoi(c.Param("statusid"))
// 	if errsid != nil {
// 		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
// 	} else {
//
// 		//check permisson
// 		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
// 		if allowed, _ := services.CheckPInteractivePermission(statusid, userid); allowed == false {
// 			libs.ResponseForbiddenJSON(c, 220, "Status not visible")
// 			return
// 		}
// 		// binding
// 		json := struct {
// 			Message string `json:"message" valid:"nonzero"`
// 		}{}
// 		if errBind := c.Bind(&json); errBind != nil {
// 			libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
// 			return
// 		}
//
// 		// validation
// 		if len(json.Message) == 0 {
// 			libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
// 			return
// 		}
//
// 		commentid, errcid := services.CreateStatusComment(statusid, userid, json.Message)
// 		if errcid == nil && commentid >= 0 {
// 			libs.ResponseCreatedJSON(c, 1, "Create status comment successful", map[string]interface{}{"id": commentid})
//
// 			// auto Increase Comments
// 			go func() {
// 				ok, err := services.IncreaseObjectComments(statusid)
// 				if err != nil {
// 					fmt.Printf("ERROR in IncreaseStatusComments service: %s", err.Error())
// 				}
// 				if ok != true {
// 					fmt.Printf("ERROR in IncreaseStatusComments service")
// 				}
// 			}()
//
// 			// push noti
// 			go func() {
// 				userComment, _ := services.GetUser(userid)
// 				status, _ := services.GetUserStatus(statusid, userid)
// 				PushTest(status.UserID, 1, "@"+userComment.Username+" bình luận bài viết của bạn", json.Message)
//
// 			}()
// 			return
// 		}
//
// 		libs.ResponseServerErrorJSON(c)
// 		if errcid != nil {
// 			fmt.Printf("ERROR in CreateStatusComment services: %s", errcid.Error())
// 		} else {
// 			fmt.Printf("ERROR in CreateStatusComment services: Don't CreateStatusComment")
// 		}
//
// 	}
// }
//
// // GetPhotoComments func
// func GetPhotoComments(c *gin.Context) {
// 	photoid, errpid := strconv.Atoi(c.Param("photoid"))
// 	if errpid != nil {
// 		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
// 	} else {
//
// 		//check permisson
// 		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
// 		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, userid); allowed == false {
// 			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
// 			return
// 		}
//
// 		sort := c.DefaultQuery("sort", "-created_at")
// 		print(sort)
// 		orderby, errSort := libs.ConvertSort(sort)
// 		if errSort != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSort.Error())
// 			return
// 		}
// 		skip, errSkip := strconv.Atoi(c.DefaultQuery("skip", "0"))
// 		if errSkip != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errSkip.Error())
// 			return
// 		}
// 		limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "25"))
// 		if errLimit != nil {
// 			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: "+errLimit.Error())
// 			return
// 		}
//
// 		commentList, errList := services.GetPhotoComments(photoid, orderby, skip, limit)
// 		if errList == nil {
// 			libs.ResponseEntityListJSON(c, 1, "Photo Comments List", commentList, nil, len(commentList))
// 			return
// 		}
//
// 		libs.ResponseServerErrorJSON(c)
// 		if errList != nil {
// 			fmt.Printf("ERROR in GetPhotoComments services: %s", errList.Error())
// 		} else {
// 			fmt.Printf("ERROR in GetPhotoComments services: Don't GetPhotoComments")
// 		}
//
// 	}
// }
//
// // CreatePhotoComment func
// func CreatePhotoComment(c *gin.Context) {
// 	photoid, errpid := strconv.Atoi(c.Param("photoid"))
// 	if errpid != nil {
// 		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
// 	} else {
//
// 		//check permisson
// 		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
// 		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, userid); allowed == false {
// 			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
// 			return
// 		}
//
// 		// binding
// 		json := struct {
// 			Message string `json:"message" valid:"nonzero"`
// 		}{}
// 		if errBind := c.Bind(&json); errBind != nil {
// 			libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
// 			return
// 		}
//
// 		// validation
// 		if len(json.Message) == 0 {
// 			libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Message is NULL", nil)
// 			return
// 		}
//
// 		commentid, errcid := services.CreatePhotoComment(photoid, userid, json.Message)
// 		if errcid == nil && commentid >= 0 {
// 			libs.ResponseCreatedJSON(c, 1, "Create status comment successful", map[string]interface{}{"id": commentid})
//
// 			// auto Increase Comments
// 			go func() {
// 				ok, err := services.IncreaseObjectComments(photoid)
// 				if err != nil {
// 					fmt.Printf("ERROR in IncreaseStatusComments service: %s", err.Error())
// 				}
// 				if ok != true {
// 					fmt.Printf("ERROR in IncreaseStatusComments service")
// 				}
// 			}()
//
// 			return
// 		}
//
// 		libs.ResponseServerErrorJSON(c)
// 		if errcid != nil {
// 			fmt.Printf("ERROR in CreateStatusComment services: %s", errcid.Error())
// 		} else {
// 			fmt.Printf("ERROR in CreateStatusComment services: Don't CreateStatusComment")
// 		}
//
// 	}
// }
