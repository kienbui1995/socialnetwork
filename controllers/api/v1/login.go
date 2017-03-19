package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/middlewares"
	"github.com/kienbui1995/socialnetwork/models"
	"github.com/kienbui1995/socialnetwork/services"
)

// SuperSecretPassword var to Sign for Token
var SuperSecretPassword = []byte("socialnetworkTLSEN")

// Login func is controller login
func Login(c *gin.Context) {
	var json models.Login
	err := c.Bind(&json)
	if err != nil {
		c.Abort()
		return
	}
	defaultvalue := ""
	if json.Username == defaultvalue || json.Password == defaultvalue || json.Device == defaultvalue {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Missing a few parameters",
		})
		c.Abort()
		return
	}
	id, err := services.Login(json)

	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		c.Abort()
		return
	}

	tokenstring, errtoken := middlewares.GenerateToken(id, json.Device, SuperSecretPassword)
	if errtoken != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": errtoken.Error(),
		})
		return
	}
	services.SaveToken(id, tokenstring)
	//c.Header("token", tokenstring)
	//tokenstruct truct
	type TokenStruct struct {
		Token string `json:"token"`
	}
	token1 := TokenStruct{Token: tokenstring}
	c.JSON(200, gin.H{
		"code":    1,
		"message": "Login successful!",
		"data":    token1,
	})
}

// Logout func to remove token of user
func Logout(c *gin.Context) {
	var json struct{ Token string }
	err := c.Bind(&json)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
	}
	// delete token from DB
	claims, err := middlewares.ExtractClaims(json.Token, SuperSecretPassword)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}
	userid := claims["userid"].(float64)
	if c.DefaultPostForm("device", "") != claims["device"].(string) {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Device is wrong.",
		})
		return
	}

	existtoken, errtoken := services.CheckExistToken(int(userid), c.PostForm("token"))
	if errtoken != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": errtoken.Error(),
		})
		return
	}
	if existtoken == true {
		deletetoken, errdeletetoken := services.DeleteToken(int(userid))
		if errdeletetoken != nil {
			c.JSON(200, gin.H{
				"code":    -1,
				"message": errtoken.Error(),
			})
			return
		}
		if deletetoken == true {
			c.JSON(200, gin.H{
				"code":    1,
				"message": "Logout successful",
			})
			return
		}
	}

}
