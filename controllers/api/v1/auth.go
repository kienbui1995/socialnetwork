package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/middlewares"
	"github.com/kienbui1995/socialnetwork/services"
)

// AuthHandler func to check user
func AuthHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if len(token) == 0 {
		err := errors.New("No token")
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	ok, err := middlewares.ValidateToken(token, SuperSecretPassword)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": err.Error(),
		})
		c.Abort()
		return
	}
	if ok == false {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "Token invalid",
		})
		c.Abort()
		return
	}
	claims, errclaim := middlewares.ExtractClaims(token, SuperSecretPassword)
	if errclaim != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": errclaim.Error(),
		})
		c.Abort()
		return
	}
	existtoken, errexist := services.CheckExistToken(int(claims["userid"].(float64)), token)
	if errexist != nil {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": errexist.Error(),
		})
		c.Abort()
		return
	}
	if existtoken != true {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "No exist token",
		})
		c.Abort()
		return
	}
}
