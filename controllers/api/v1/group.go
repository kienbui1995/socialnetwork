package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// CreateGroup func  to create a group
func CreateGroup(c *gin.Context) {

	//check permisson
	userid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
	if userid < 0 || errGet != nil {
		libs.ResponseAuthJSON(c, 200, "Permissions error")
		return
	}

	json := models.Group{}
	if errBind := c.Bind(&json); errBind != nil {
		libs.ResponseJSON(c, 400, 100, "Invalid parameter: "+errBind.Error(), nil)
		return
	}

	// validation
	if len(json.Name) == 0 {
		libs.ResponseJSON(c, 400, 100, "Missing a few fields:  Name is NULL", nil)
		return
	}
	if json.Status == 0 {
		json.Status = 1
	}
	if json.Privacy == 0 {
		json.Privacy = 1
	}

	groupid, errgid := services.CreateGroup(userid, json)
	if errgid == nil && groupid >= 0 {
		libs.ResponseSuccessJSON(c, 1, "Create user post successful", map[string]interface{}{"id": groupid})

		return

		// auto Increase Posts
		// go func() {
		// 	ok, err := services.IncreasePosts(userid)
		// 	if err != nil {
		// 		fmt.Printf("ERROR in IncreasePosts service: %s", err.Error())
		// 	}
		// 	if ok != true {
		// 		fmt.Printf("ERROR in IncreasePosts service")
		// 	}
		// }()

		// push noti
		// go func() {
		// 	user, _ := services.GetUser(userid)
		// 	ids, errGetIDs := services.GetFollowers(userid)
		// 	if len(ids) > 0 && errGetIDs == nil {
		// 		for index := 0; index < len(ids); index++ {
		// 			PushTest(ids[index].UserID, postID, "post", "@"+user.Username+action, json.Message)
		// 		}
		//
		// 	}
		// }()

	}
	libs.ResponseServerErrorJSON(c)
	if errgid != nil {
		fmt.Printf("ERROR in CreateGroup services: %s", errgid.Error())
	} else {
		fmt.Printf("ERROR in CreateGroup services: Don't create User Group")
	}

}

// GetGroups func to list of group
func GetGroups(c *gin.Context) {
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
		typePost := c.Query("type")
		ItypePost := 0
		if typePost != "class" && typePost != "all" && len(typePost) > 0 {
			libs.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: type")
			return
		} else if typePost == "class" {
			ItypePost = configs.PostPhoto
		} else if typePost == "all" {
			ItypePost = configs.PostStatus
		}
		sort := c.DefaultQuery("sort", "-created_at")
		print(sort)
		orderby, errSort := libs.ConvertSort(sort)
		if errSort != nil {
			libs.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errSort.Error())
			return
		}
		skip, errSkip := strconv.Atoi(c.DefaultQuery("skip", "0"))
		if errSkip != nil {
			libs.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errSkip.Error())
			return
		}
		limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "25"))
		if errLimit != nil {
			libs.ResponseBadRequestJSON(c, configs.EcParam, "Invalid parameter: "+errLimit.Error())
			return
		}

		postList, errList := services.GetUserPosts(userid, myuserid, orderby, skip, limit, ItypePost)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, "User Post List", postList, nil, len(postList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetUserPosts services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetUserPosts services: Don't get User Posts")
		}

	}
}

// GetGroup func to get info of the group
func GetGroup(c *gin.Context) {
	groupid, errgid := strconv.Atoi(c.Param("groupid"))
	if errgid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: groupid")
	} else {

		//check permisson
		myuserid, errID := GetUserIDFromToken(c.Request.Header.Get("token"))
		if errID != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		group, errGet := services.GetGroup(groupid, myuserid)
		if errGet == nil {
			libs.ResponseSuccessJSON(c, 1, "get a Group", group)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errGet != nil {
			fmt.Printf("ERROR in GetUserPosts services: %s", errGet.Error())
		} else {
			fmt.Printf("ERROR in GetUserPosts services: Don't get User Posts")
		}

	}
}

// UpdateGroup func to update the post ~doing
func UpdateGroup(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid post id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPost(postid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDByPostID(postid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		json := struct {
			Message string `json:"message"`
			Photo   string `json:"photo"`
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

		post, errUpdate := services.UpdateUserPost(postid, json.Message, json.Photo, json.Privacy, json.Status)
		if errUpdate == nil && post.CreatedAt > 0 {
			libs.ResponseSuccessJSON(c, 1, "Update user post successful", post)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errUpdate != nil {
			fmt.Printf("ERROR in UpdateUserPost services: %s", errUpdate.Error())
		} else {
			fmt.Printf("ERROR in UpdateUserPost services: Don't update User Posts")
		}

	}
}

// DeleteGroup func to delete a group ~doing
func DeleteGroup(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid post id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPost(postid); exist != true {
			libs.ResponseBadRequestJSON(c, 2, "No exist this object")
			return
		}

		userid, _ := services.GetUserIDByPostID(postid)
		//check permisson
		if id, errGet := GetUserIDFromToken(c.Request.Header.Get("token")); userid != id || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		ok, errDel := services.DeleteUserPost(postid)
		if errDel == nil && ok == true {
			libs.ResponseSuccessJSON(c, 1, "Delete user post successful", nil)

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
			fmt.Printf("ERROR in DeleteUserPost services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeleteUserPost services: Don't delete User Posts")
		}

	}
}

// GetGroupFeed func to feed on group ~doing
func GetGroupFeed(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid post id")

	} else {

		//check exist
		if exist, _ := services.CheckExistUserPost(postid); exist != true {
			libs.ResponseNotFoundJSON(c, 2, "No exist this object")
			return
		}

		// userid, _ := services.GetUserIDPostedStatus(postid)
		//check permisson ~needfix when privacy not public
		myuserid, errGet := GetUserIDFromToken(c.Request.Header.Get("token"))
		if myuserid == -1 || errGet != nil {
			libs.ResponseAuthJSON(c, 200, "Permissions error")
			return
		}

		post, errGet := services.GetUserPost(postid, myuserid)
		if errGet == nil && post.PostID == postid {
			libs.ResponseSuccessJSON(c, 1, "Get user post successful", post)
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errGet != nil {
			fmt.Printf("ERROR in GetUserPost services: %s", errGet.Error())
		} else {
			fmt.Printf("ERROR in GetUserPost services: Don't get User post")
		}

	}
}

// CreateGroupPost func ~doing
func CreateGroupPost(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// check liked
		if liked, _ := services.CheckExistPostLike(postid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		likes, errLike := services.CreatePostLike(postid, userid)
		if errLike == nil && likes >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

			// auto Increase post Likes
			go func() {
				ok, err := services.IncreasePostLikes(postid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePostLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePostLikes service")
				}
			}()

			// push noti
			go func() {
				post, _ := services.GetUserPost(postid, userid)
				userLiked, _ := services.GetUser(userid)
				PushTest(post.UserID, postid, "post", "@"+userLiked.Username+" thích trạng thái của bạn", post.Message)

			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreatePostLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreatePostLike services: Don't CreatePostLike")
		}

	}
}

// CreateJoinGroup func ~doing
func CreateJoinGroup(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// check liked
		if liked, _ := services.CheckExistPostLike(postid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		likes, errLike := services.CreatePostLike(postid, userid)
		if errLike == nil && likes >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

			// auto Increase post Likes
			go func() {
				ok, err := services.IncreasePostLikes(postid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePostLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePostLikes service")
				}
			}()

			// push noti
			go func() {
				post, _ := services.GetUserPost(postid, userid)
				userLiked, _ := services.GetUser(userid)
				PushTest(post.UserID, postid, "post", "@"+userLiked.Username+" thích trạng thái của bạn", post.Message)

			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreatePostLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreatePostLike services: Don't CreatePostLike")
		}

	}
}

// CreateJoinGroupRequest func ~doing
func CreateJoinGroupRequest(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// check liked
		if liked, _ := services.CheckExistPostLike(postid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		likes, errLike := services.CreatePostLike(postid, userid)
		if errLike == nil && likes >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

			// auto Increase post Likes
			go func() {
				ok, err := services.IncreasePostLikes(postid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePostLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePostLikes service")
				}
			}()

			// push noti
			go func() {
				post, _ := services.GetUserPost(postid, userid)
				userLiked, _ := services.GetUser(userid)
				PushTest(post.UserID, postid, "post", "@"+userLiked.Username+" thích trạng thái của bạn", post.Message)

			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreatePostLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreatePostLike services: Don't CreatePostLike")
		}

	}
}

// GetGroupMembers func ~doing
func GetGroupMembers(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
	} else {

		//check permisson
		userid, _ := GetUserIDFromToken(c.Request.Header.Get("token"))

		// check liked
		if liked, _ := services.CheckExistPostLike(postid, userid); liked == true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		likes, errLike := services.CreatePostLike(postid, userid)
		if errLike == nil && likes >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Like post successful", map[string]int{"likes": likes})

			// auto Increase post Likes
			go func() {
				ok, err := services.IncreasePostLikes(postid)
				if err != nil {
					fmt.Printf("ERROR in IncreasePostLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in IncreasePostLikes service")
				}
			}()

			// push noti
			go func() {
				post, _ := services.GetUserPost(postid, userid)
				userLiked, _ := services.GetUser(userid)
				PushTest(post.UserID, postid, "post", "@"+userLiked.Username+" thích trạng thái của bạn", post.Message)

			}()
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errLike != nil {
			fmt.Printf("ERROR in CreatePostLike services: %s", errLike.Error())
		} else {
			fmt.Printf("ERROR in CreatePostLike services: Don't CreatePostLike")
		}

	}
}
