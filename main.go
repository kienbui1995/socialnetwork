package main

import (
	"github.com/gin-gonic/gin"
	apiv1 "github.com/kienbui1995/socialnetwork/controllers/api/v1"
)

var (
	neo4jURL = "Bolt://neo4j:tlis2016@localhost:7687"
)

func main() {
	router := gin.Default()
	// Work for login
	router.POST("/login", apiv1.Login)     // login method
	router.POST("/logout", apiv1.Logout)   // logout method
	router.POST("/user", apiv1.CreateUser) // create a user
	// Work for User
	RUser := router.Group("/user", apiv1.AuthHandler)
	{
		// //Check token in Header
		// RUser.Use(apiv1.AuthHandler)

		// user
		RUser.GET("", apiv1.GetUser)            // get a few user
		RUser.GET("/:userid", apiv1.GetUser)    // get a user
		RUser.PUT("/:userid", apiv1.UpdateUser) // update a user

		RUser.DELETE("/:userid", apiv1.DeleteUser) // delete a user

		// user with post
		RUser.POST("/:userid/post", apiv1.CreatePost)           // user create a post
		RUser.GET("/:userid/post", apiv1.GetPost)               // get post of this post
		RUser.GET("/:userid/post/:postid", apiv1.GetPost)       // user get a own post
		RUser.PUT("/:userid/post/:postid", apiv1.UpdatePost)    // user update a own post
		RUser.DELETE("/:userid/post/:postid", apiv1.DeletePost) // user delete a own post

		//user with login
		// /RUser.POST("/login", apiv1.Login)
	}

	// Work for Post
	RPost := router.Group("/post")
	{
		RPost.GET("", apiv1.GetPost)              // get a few post
		RPost.GET("/:postid", apiv1.GetPost)      // get a post
		RPost.PUT("/:postid", apiv1.UpdatePost)   // update a post
		RPost.POST("", apiv1.CreatePost)          // create a post
		RPost.DELETE(":postid", apiv1.DeletePost) // delete a post
	}

	router.Run(":8080")

}
