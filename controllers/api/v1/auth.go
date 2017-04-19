package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/kienbui1995/socialnetwork/middlewares"
	"github.com/kienbui1995/socialnetwork/services"
)

// AuthHandler func to check user
func AuthHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if len(token) == 0 {
		libs.ResponseAuthJSON(c, 404, "Missing token.")
		return
	}
	ok, err := middlewares.ValidateToken(token, SuperSecretPassword)
	if err != nil {
		libs.ResponseAuthJSON(c, 407, "Error in checking toke: "+err.Error())
		return
	}
	if ok == false {
		libs.ResponseAuthJSON(c, 405, "Invalid token.")
		return
	}
	claims, errclaim := middlewares.ExtractClaims(token, SuperSecretPassword)
	if errclaim != nil {
		libs.ResponseAuthJSON(c, 407, "Error in extracting claims in token: "+errclaim.Error())
		return
	}
	existtoken, errexist := services.CheckExistToken(int(claims["userid"].(float64)), token)
	if errexist != nil {
		libs.ResponseAuthJSON(c, 407, "Error in checking token: "+errexist.Error())
		return
	}
	if existtoken != true {
		libs.ResponseAuthJSON(c, 406, "No exist token.")
		return
	}
}
