package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/services"
)

// GetStatusComments func
func GetStatusComments(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
	} else {

		//check permisson

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

		commentList, errList := services.GetStatusComments(statusid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Statuses Comments List", commentList, nil, len(commentList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetStatusComments services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetStatusComments services: Don't getGetStatusComments")
		}

	}
}

// CreateStatusComment func
func CreateStatusComment(c *gin.Context) {
	statusid, errsid := strconv.Atoi(c.Param("statusid"))
	if errsid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errsid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// binding
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

		commentid, errcid := services.CreateStatusComment(statusid, userid, json.Message)
		if errcid == nil && commentid >= 0 {
			libs.ResponseCreatedJSON(c, 1, "Create status comment successful", map[string]interface{}{"id": commentid})
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errcid != nil {
			fmt.Printf("ERROR in CreateStatusComment services: %s", errcid.Error())
		} else {
			fmt.Printf("ERROR in CreateStatusComment services: Don't CreateStatusComment")
		}

	}
}

// UpdateStatusComment func
func UpdateStatusComment(c *gin.Context) {
	commentid, errcid := strconv.Atoi(c.Param("commentid"))
	if errcid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserStatus(commentid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
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

		updated, errUpdate := services.UpdateStatusComment(commentid, json.Message)
		if errUpdate == nil && updated == true {
			libs.ResponseSuccessJSON(c, 1, "Update comment successful", nil)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errUpdate != nil {
			fmt.Printf("ERROR in UpdateStatusComment services: %s", errUpdate.Error())
		} else {
			fmt.Printf("ERROR in UpdateStatusComment services: Don't UpdateStatusComment")
		}
	}
}

// DeleteStatusComment func to delete a comment
func DeleteStatusComment(c *gin.Context) {
	commentid, errcid := strconv.Atoi(c.Param("commentid"))
	if errcid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistComment(commentid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDWroteComment(commentid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		ok, errDel := services.DeleteStatusComment(commentid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Delete comment successful", nil)
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
