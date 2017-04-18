package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/libs"
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
		libs.ResponseAuthJSON(c, -1, "Missing a few parameters")
		c.Abort()
		return
	}
	id, err := services.Login(json)

	if err != nil {
		libs.ResponseAuthJSON(c, -1, err.Error())
		c.Abort()
		return
	}
	tokenstring, errtoken := middlewares.GenerateToken(id, json.Device, SuperSecretPassword)
	if errtoken != nil {
		libs.ResponseAuthJSON(c, -1, errtoken.Error())
		return
	}
	go services.SaveToken(id, json.Device, tokenstring)
	token := map[string]string{"token": tokenstring}
	libs.ResponseSuccessJSON(c, 1, "Login successful!", token)
}

// Logout func to remove token of user
func Logout(c *gin.Context) {

	token := c.Request.Header.Get("token")

	// delete token from DB
	claims, err := middlewares.ExtractClaims(token, SuperSecretPassword)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		return
	}

	userid := claims["userid"].(float64)
	existToken, errExistToken := services.CheckExistToken(int(userid), token)
	if errExistToken != nil || existToken == false {
		libs.ResponseAuthJSON(c, -1, errExistToken.Error())
		return
	}
	deletetoken, errdelete := services.DeleteToken(int(userid), token)
	if deletetoken != true || errdelete != nil {
		libs.ResponseBadRequestJSON(c, -1, errdelete.Error())
		return
	}
	libs.ResponseNoContentJSON(c, 1, "Logout successful")

}
