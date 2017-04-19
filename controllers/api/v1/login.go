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
		libs.ResponseAuthJSON(c, 101, "Missing a few fields.")
		c.Abort()
		return
	}
	id, err := services.Login(json)

	if err != nil {
		libs.ResponseAuthJSON(c, 409, "No exist user: "+err.Error())
		c.Abort()
		return
	}
	tokenstring, errtoken := middlewares.GenerateToken(id, json.Device, SuperSecretPassword)
	if errtoken != nil {
		libs.ResponseAuthJSON(c, 408, "Error in generate token: "+errtoken.Error())
		return
	}
	if saveToken, err := services.SaveToken(id, json.Device, tokenstring); saveToken != true || err != nil {
		libs.ResponseBadRequestJSON(c, 1, "Don't save token"+err.Error())
		return
	}
	token := map[string]string{"token": tokenstring}
	libs.ResponseSuccessJSON(c, 1, "Login successful!", token)
}

// LoginViaFacebook func is login or sign up via Facebook
func LoginViaFacebook(c *gin.Context) {
	type FacebookToken struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
	}
	var json FacebookToken
	err := c.Bind(&json)
	if err != nil {
		c.Abort()
		return
	}

	verify := libs.VerifyFacebookID(json.ID, json.AccessToken)
	if verify != true {

	}
}

// Logout func to remove token of user
func Logout(c *gin.Context) {

	token := c.Request.Header.Get("token")

	// delete token from DB
	claims, err := middlewares.ExtractClaims(token, SuperSecretPassword)
	if err != nil {
		libs.ResponseBadRequestJSON(c, 407, "Error in checking token: "+err.Error())
		// c.JSON(200, gin.H{
		// 	"code":    -1,
		// 	"message": err.Error(),
		// })
		return
	}

	userid := claims["userid"].(float64)
	existToken, errExistToken := services.CheckExistToken(int(userid), token)
	if errExistToken != nil || existToken == false {
		libs.ResponseAuthJSON(c, 406, "No exist token: "+errExistToken.Error())
		return
	}
	deletetoken, errdelete := services.DeleteToken(int(userid), token)
	if deletetoken != true || errdelete != nil {
		libs.ResponseBadRequestJSON(c, 407, "Error in checking token: "+errdelete.Error())
		return
	}
	libs.ResponseNoContentJSON(c, 1, "Logout successful")
}
