package v1

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/services"
)

// CreateUserPost func  to create a user status/photo upload
func CreateUserPost(c *gin.Context) {
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
			Photo   string `json:"photo"`
			Message string `json:"message"`
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
		if json.Status == 0 {
			json.Status = 1
		}
		if json.Privacy == 0 {
			json.Privacy = 1
		}
		action := " cập nhật trạng thái"
		if len(json.Photo) > 0 {
			action = " đăng ảnh"
		}
		postID, errpid := services.CreateUserPost(userid, json.Message, json.Photo, json.Privacy, json.Status)
		if errpid == nil && postID >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Create user post successful", map[string]interface{}{"id": postID})

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
				user, _ := services.GetUser(userid)
				ids, errGetIDs := services.GetFollowers(userid)
				if len(ids) > 0 && errGetIDs == nil {
					for index := 0; index < len(ids); index++ {
						PushTest(ids[index].UserID, 1, "@"+user.Username+action, json.Message)
					}

				}
			}()

			return
		}
		libs.ResponseServerErrorJSON(c)
		if errpid != nil {
			fmt.Printf("ERROR in CreateUserPost services: %s", errpid.Error())
		} else {
			fmt.Printf("ERROR in CreateUserPost services: Don't create User Post")
		}
	}
}

// GetUserPosts func to list of post
func GetUserPosts(c *gin.Context) {
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
		if typePost != "photo" && typePost != "status" && len(typePost) > 0 {
			libs.ResponseBadRequestJSON(c, configs.APIEcParam, "Invalid parameter: type")
			return
		} else if typePost == "photo" {
			ItypePost = configs.PostPhoto
		} else if typePost == "status" {
			ItypePost = configs.PostStatus
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

// UpdateUserPost func to update the post
func UpdateUserPost(c *gin.Context) {
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

// DeleteUserPost func to delete a user post
func DeleteUserPost(c *gin.Context) {
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

// GetUserPost func to delete a user post
func GetUserPost(c *gin.Context) {
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

// CreatePostLike func
func CreatePostLike(c *gin.Context) {
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
				PushTest(post.UserID, 1, "@"+userLiked.Username+" thích trạng thái của bạn", post.Message)

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

// DeletePostLike func to delete a like
func DeletePostLike(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid comment id")
	} else {

		//check exist
		if exist, _ := services.CheckExistUserPost(postid); exist != true {
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
		if liked, _ := services.CheckExistPostLike(postid, userid); liked != true {
			libs.ResponseBadRequestJSON(c, 3, "Exist this object: Likes")
			return
		}

		likes, errDel := services.DeletePostLike(postid, userid)
		if errDel == nil && likes >= 0 {
			libs.ResponseSuccessJSON(c, 1, "Unlike successful", map[string]int{"likes": likes})

			// auto Decrease post Likes
			go func() {
				ok, err := services.DecreasePostLikes(postid)
				if err != nil {
					fmt.Printf("ERROR in DecreasePostLikes service: %s", err.Error())
				}
				if ok != true {
					fmt.Printf("ERROR in DecreasePostLikes service")
				}
			}()

			return
		}

		libs.ResponseServerErrorJSON(c)
		if errDel != nil {
			fmt.Printf("ERROR in DeletePostLike services: %s", errDel.Error())
		} else {
			fmt.Printf("ERROR in DeletePostLike services: Don't DeletePostLike")
		}

	}
}

// GetPostLikes func
func GetPostLikes(c *gin.Context) {
	postid, errpid := strconv.Atoi(c.Param("postid"))
	if errpid != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter: "+errpid.Error())
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

		likeList, errList := services.GetPostLikes(postid, userid, orderby, skip, limit)
		if errList == nil {
			libs.ResponseEntityListJSON(c, 1, " Posts Likes User List", likeList, nil, len(likeList))
			return
		}

		libs.ResponseServerErrorJSON(c)
		if errList != nil {
			fmt.Printf("ERROR in GetPostLikes services: %s", errList.Error())
		} else {
			fmt.Printf("ERROR in GetPostLikes services: Don't GetPostLikes")
		}

	}
}
