package main

import (
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	apiv1 "github.com/kienbui1995/socialnetwork/controllers/api/v1"
)

func main() {
	myfile, _ := os.Create("server.log")

	router := gin.Default()
	router.Use(io.MultiWriter(myfile, os.Stdout))
	// Func Test
	router.POST("/test/users", apiv1.CreateUserTest)
	// router.POST("/test/push", apiv1.PushTest)
	// Work for login
	router.POST("/login", apiv1.Login)                             // login method1
	router.POST("/login_facebook", apiv1.LoginViaFacebook)         // login via facebook2
	router.POST("/logout", apiv1.Logout)                           // logout method3
	router.POST("/sign_up", apiv1.SignUp)                          // create a user4
	router.POST("/forgot_password", apiv1.ForgotPassword)          // forgot password method5
	router.POST("/verify_recovery_code", apiv1.VerifyRecoveryCode) // verify recovery code method6
	router.PUT("/renew_password", apiv1.RenewPassword)             // renew password after verify_recovery_code7

	router.GET("/ws", apiv1.WsHandler)

	// Work for User
	authorized := router.Group("/", apiv1.AuthHandler)
	{
		authorized.GET("/find_user", apiv1.FindUser) // find user by username or fullname 8
		RUser := authorized.Group("/users")
		{
			// user
			RUser.GET("", apiv1.GetUser)            // get a few user ~needfix9
			RUser.GET("/:userid", apiv1.GetUser)    // get a user  ~needfix10
			RUser.PUT("/:userid", apiv1.UpdateUser) // update a user ~needfix11

			RUser.DELETE("/:userid", apiv1.DeleteUser) // delete a user12

			// user with post
			RUser.POST("/:userid/posts", apiv1.CreateUserPost)     // create a post on the user's wall via userid13
			RUser.GET("/:userid/posts", apiv1.GetUserPosts)        // get a posts list on the user's wall via userid14
			RUser.GET("/:userid/posts/:postid", apiv1.GetUserPost) // user get a own post15

			// user witth follow
			RUser.POST("/:userid/subscribers/:toid", apiv1.CreateUserSubscribers)   // create follow from userid to toid16
			RUser.DELETE("/:userid/subscribers/:toid", apiv1.DeleteUserSubscribers) // unfollow follow from userid to toid17
			RUser.GET("/:userid/followers", apiv1.GetFollowers)                     // get users list who being follow userid18
			RUser.GET("/:userid/subscribers", apiv1.GetSubscribers)                 // get users list whose userid being follow19

			// user with status
			// RUser.POST("/:userid/statuses", apiv1.CreateUserStatus)            // create a status on the user's wall via userid20
			// RUser.GET("/:userid/statuses", apiv1.GetUserStatuses)              // get a statuses list on the user's wall via userid21
			// RUser.PUT(":userid/statuses/:statusid", apiv1.UpdateUserStatus)    // update a user status via statusid
			// RUser.DELETE(":userid/statuses/:statusid", apiv1.DeleteUserStatus) // delete a user status via statusid

			// user with post
			// RUser.POST("/:userid/photos", apiv1.CreateUserPhoto)      // create a photo on the user's wall via userid13
			// RUser.GET("/:userid/photos", apiv1.GetUserPhotos)         // get a photos list on the user's wall via userid14
			// RUser.GET("/:userid/photos/:photoid", apiv1.GetUserPhoto) // user get a own photo15
			// RUser.DELETE("/:userid/photos/:photoid", apiv1.DeleteUserPhoto)

			RUser.GET("/:userid/home", apiv1.GetNewsFeed) // get newsfeed of user by userid22
			// RUser.GET("/:userid/feed", apiv1.GetUserWall)           // get post and status on the user's wall via userid23
			// RUser.POST("/:userid/posts", apiv1.CreateWallPost) // create a post or a status on the user's wall via userid24
			// RUser.GET("/:userid/groups", apiv1.GetUserJoinedGroups) // get a groups list that user joined via userid25
		}

		// Work for Post
		RPost := authorized.Group("/posts")
		{
			RPost.GET("", apiv1.GetUserPosts)             // get a posts list 26
			RPost.GET("/:postid", apiv1.GetUserPost)      // get a post 27
			RPost.PUT("/:postid", apiv1.UpdateUserPost)   // update a post 28
			RPost.POST("", apiv1.CreateUserPost)          // create a post 29
			RPost.DELETE(":postid", apiv1.DeleteUserPost) // delete a post 30

			// post with comment
			RPost.GET("/:postid/comments", apiv1.GetComments)    // get a comments list on the post via postid 31
			RPost.POST("/:postid/comments", apiv1.CreateComment) // create a comment on the post via postid 32

			// post with like ~needfix can react
			RPost.GET("/:postid/likes", apiv1.GetPostLikes)      // get a users list who liked post via postid 33
			RPost.POST("/:postid/likes", apiv1.CreatePostLike)   // create a like on the post via postid 34
			RPost.DELETE("/:postid/likes", apiv1.DeletePostLike) // unlike on the post via postid 35
		}

		// // Work for Status
		// RStatus := authorized.Group("/statuses")
		// {
		// RStatus.GET("/:statusid", apiv1.GetUserStatus)       // get a status via statusid 36
		// RStatus.PUT("/:statusid", apiv1.UpdateUserStatus)    // update a user status via statusid
		// RStatus.DELETE("/:statusid", apiv1.DeleteUserStatus) // delete a user status via statusid
		//
		// // status with comment
		// RStatus.GET("/:statusid/comments", apiv1.GetStatusComments)    // get a comments list on the status via statusid 37
		// RStatus.POST("/:statusid/comments", apiv1.CreateStatusComment) // create a comment on the status via statusid 38
		//
		// // Status with like ~needfix can react
		// RStatus.GET("/:statusid/likes", apiv1.GetStatusLikes)      // get a users list who liked status via statusid 39
		// RStatus.POST("/:statusid/likes", apiv1.CreateStatusLike)   // create a like on the status via statusid 40
		// RStatus.DELETE("/:statusid/likes", apiv1.DeleteStatusLike) // unlike on the status via statusid 41
		// }

		// // Work for Photo
		// RPhoto := authorized.Group("/photos")
		// {
		// 	RPhoto.GET("/:photoid", apiv1.GetUserPhoto)       // get a photo via photoid 36
		// 	RPhoto.PUT("/:photoid", apiv1.UpdateUserPhoto)    // update a user photo via photoid
		// 	RPhoto.DELETE("/:photoid", apiv1.DeleteUserPhoto) // delete a user photo via photoid
		//
		// 	// Photo with comment
		// 	RPhoto.GET("/:photoid/comments", apiv1.GetPhotoComments)    // get a comments list on the photo via photoid 37
		// 	RPhoto.POST("/:photoid/comments", apiv1.CreatePhotoComment) // create a comment on the photo via photoid 38
		//
		// 	// Photo with like ~needfix can react
		// 	RPhoto.GET("/:photoid/likes", apiv1.GetPhotoLikes)      // get a users list who liked photo via photoid 39
		// 	RPhoto.POST("/:photoid/likes", apiv1.CreatePhotoLike)   // create a like on the photo via photoid 40
		// 	RPhoto.DELETE("/:photoid/likes", apiv1.DeletePhotoLike) // unlike on the photo via photoid 41
		// }

		// Work for Comment
		RComment := authorized.Group("/comments")
		{
			// RComment.POST("", apiv1.CreateComment)
			RComment.GET("/:commentid", apiv1.GetComment) // get a comment via commentid 42
			RComment.PUT("/:commentid", apiv1.UpdateComment)
			RComment.DELETE("/:commentid", apiv1.DeleteComment) //delete a comment via commentid
			// Comment with like ~needfix can react
			// RComment.GET("/:commentid/likes", apiv1.GetCommentLikes)      // get a users list who liked stacommenttus via commentid 43
			// RComment.POST("/:commentid/likes", apiv1.CreateCommentLike)   // create a like on the comment via commentid 44
			// RComment.DELETE("/:commentid/likes", apiv1.DeleteCommentLike) // unlike on the comment via commentid 45
		}

		// Work for Group
		RGroup := authorized.Group("/groups")
		{
			//
			RGroup.POST("", apiv1.CreateGroup)            // create a group
			RGroup.GET("", apiv1.GetGroups)               // get group list
			RGroup.GET("/:groupid", apiv1.GetGroup)       // get info of a group by groupid 46
			RGroup.PUT("/:groupid", apiv1.UpdateGroup)    // update info of a group by groupid 50
			RGroup.DELETE("/:groupid", apiv1.DeleteGroup) // delete a group by groupid 51

			// group with post
			RGroup.POST("/:groupid/posts", apiv1.CreateGroupPost) // create a post or a status to a group by groupid 47
			RGroup.GET("/:groupid/posts", apiv1.GetGroupFeed)     // get a posts list in a group by groupid 48

			// group with members
			RGroup.POST(":groupid/members", apiv1.CreateJoinGroup)                 // create join group
			RGroup.POST(":groupid/members/requests", apiv1.CreateJoinGroupRequest) // create Request to Join a Private Group
			RGroup.GET(":groupid/members/requests", apiv1.GetJoinGroupRequest)     // get Request list to Join a Private Group

			// Replaced by  PUT /group-memberships-requests/requestid
			RGroup.PUT(":groupid/members/requests", apiv1.GetJoinGroupRequest) // get Request list to Join a Private Group

			RGroup.GET("/:groupid/members", apiv1.GetGroupMembers) // get a users list in a group by groupid 49
		}

	}
	router.Run(":" + strconv.Itoa(configs.APIPort))

}
