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
		libs.ResponseAuthJSON(c, -1, "No token")
		return
	}
	ok, err := middlewares.ValidateToken(token, SuperSecretPassword)
	if err != nil {
		libs.ResponseAuthJSON(c, -1, err.Error())
		return
	}
	if ok == false {
		libs.ResponseAuthJSON(c, -1, "Token invalid")
		return
	}
	claims, errclaim := middlewares.ExtractClaims(token, SuperSecretPassword)
	if errclaim != nil {
		libs.ResponseAuthJSON(c, -1, errclaim.Error())
		return
	}
	existtoken, errexist := services.CheckExistToken(int(claims["userid"].(float64)), token)
	if errexist != nil {
		libs.ResponseAuthJSON(c, -1, errexist.Error())
		return
	}
	if existtoken != true {
		libs.ResponseAuthJSON(c, -1, "No exist token")
		return
	}
}
