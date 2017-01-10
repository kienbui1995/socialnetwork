package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/middlewares"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// SuperSecretPassword var to Sign for Token
var SuperSecretPassword = []byte("socialnetworkTLSEN")

//Login func is controller login
func Login(c *gin.Context) {
	defaultvalue := "test"
	username := c.DefaultPostForm("username", defaultvalue)
	password := c.DefaultPostForm("password", defaultvalue)
	device := c.DefaultPostForm("device", defaultvalue)
	if username == defaultvalue || password == defaultvalue || device == defaultvalue {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Missing a few parameters",
		})
		c.Abort()
		return
	}
	id, err := services.Login(models.Login{Username: username, Password: password})

	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	tokenstring, errtoken := middlewares.GenerateToken(id, device, SuperSecretPassword)
	if errtoken != nil {
		c.JSON(200, gin.H{
			"code":  -1,
			"error": errtoken,
		})
		return
	}
	c.Header("token", tokenstring)
	c.JSON(200, gin.H{
		"code":    1,
		"message": "Login successful!",
	})
}
