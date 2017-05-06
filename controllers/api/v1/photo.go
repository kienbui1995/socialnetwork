package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/services"
)

// CreateUserPhoto func
func CreateUserPhoto(c *gin.Context) {
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
			Privacy int    `json:"privacy"`
			Photo   string `json:"photo"`
			Status  int    `json:"status"`
		}{}
		if errBind := c.Bind(&json); errBind != nil {
			libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errBind.Error())
			return
		}

		// validation
		if len(json.Photo) == 0 {
			libs.ResponseBadRequestJSON(c, 100, "Missing a few fields:  Photo is NULL")
			return
		}
		if json.Privacy == 0 {
			json.Privacy = 1
		}

		if json.Status == 0 {
			json.Status = 1
		}

		photoID, errpid := services.CreateUserPhoto(userid, json.Message, json.Photo, json.Privacy, json.Status)
		if errpid == nil && photoID >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Create user status successful", map[string]interface{}{"id": photoID})

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
			return
		}
		libs.ResponseServerErrorJSON(c)
		if errpid != nil {
			fmt.Printf("ERROR in CreateUserPhoto services: %s", errpid.Error())
		} else {
			fmt.Printf("ERROR in CreateUserPhoto services: Don't CreateUserPhoto")
		}
	}
}

// GetUserPhotos func to get  a  photo
func GetUserPhotos(c *gin.Context) {
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

		photoList, errList := services.GetUserPhotos(userid, myuserid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Photos List", photoList, nil, len(photoList))
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

// UpdateUserPhoto func to  update for a photo via photoid
func UpdateUserPhoto(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid photo id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPhoto(photoid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDUploadedPhoto(photoid)
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

		status, errUpdate := services.UpdateUserPhoto(photoid, json.Message, json.Privacy, json.Status)
		if errUpdate == nil && status.UpdatedAt > 0 {
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

// DeleteUserPhoto func to delete a user photo via photoid
func DeleteUserPhoto(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid photo id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPhoto(photoid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDUploadedPhoto(photoid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		ok, errDel := services.DeleteUserPhoto(photoid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Delete user photo successful", nil)

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
			fmt.Printf("ERROR in DeleteUserPhoto services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeleteUserPhoto services: Don't delete user photo")
		}

	}
}

// GetUserPhoto func to get a user photo via photoid
func GetUserPhoto(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid photo id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPhoto(photoid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		// userid, _ := services.GetUserIDPostedStatus(statusid)
		//check permisson
		myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if myuserid == -1 || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}
		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, myuserid); allowed == false {
			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
			return
		}

		photo, errGet := services.GetUserPhoto(photoid, myuserid)
		if errGet == nil && photo.PhotoID == photoid {
			libs.ResponseSuccessJSON(c, 1, "Get user photo successful", photo)
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

// CreatePhotoLike func
func CreatePhotoLike(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, userid); allowed == false {
			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
			return
		}

		// check liked
		if liked, _ := services.CheckExistPhotoLike(photoid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		liked, errLike := services.CreatePhotoLike(photoid, userid)
		if errLike == nil && liked == true {
			libs.ResponseSuccessJSON(c, 1, "Like photo successful", nil)

			// auto Increase Status Likes
			go func() {
				ok, err := services.IncreasePhotoLikes(photoid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePhotoLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePhotoLikes service")
				}
			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreatePhotoLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreatePhotoLike services: Don't CreatePhotoLike")
		}

	}
}

// DeletePhotoLike func to delete a like
func DeletePhotoLike(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid photo id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPhoto(photoid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		//check permisson
		userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}
		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, userid); allowed == false {
			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
			return
		}

		// check liked
		if liked, _ := services.CheckExistPhotoLike(photoid, userid); liked != true {
			libs.ResponseNotFoundJSON(c, 2, "Don't Exist this object: Likes")
			return
		}

		ok, errDel := services.DeletePhotoLike(photoid, userid)
		if errDel == nil && ok == true {
			libs.ResponseJSON(c, 200, 1, "Unlike successful", nil)

			// auto Decrease Status Likes
			go func() {
				ok, err := services.DecreasePhotoLikes(photoid)
				if err != nil {
					fmt.Printf("ERROR in DecreasePhotoLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in DecreasePhotoLikes service")
				}
			}()

			return
		}

		libs.ResponseServerErrorJSON(c)
		if errDel != nil {
			fmt.Printf("ERROR in DeletePhotoLike services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeletePhotoLike services: Don't DeletePhotoLike")
		}

	}
}

// GetPhotoLikes func
func GetPhotoLikes(c *gin.Context) {
	photoid, errpid := strconv.Atoi(c.Param("photoid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))
		if allowed, _ := services.CheckPhotoInteractivePermission(photoid, userid); allowed == false {
			libs.ResponseForbiddenJSON(c, 221, "Photo not visible")
			return
		}
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

		likeList, errList := services.GetPhotoLikes(photoid, userid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, " Photo Likes User List", likeList, nil, len(likeList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetPhotoLikes services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetPhotoLikes services: Don't GetPhotoLikes")
		}

	}
}
