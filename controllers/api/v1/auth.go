package v1

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/middlewares"
)

// AuthHandler func to check user
func AuthHandler(c *gin.Context) {
	token := c.Request.Header.Get("token")
	//fmt.Printf("\n%s, %d\n", token, len(token))

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
			"message": "Token invalid!",
		})
		c.Abort()
		return
	}

}
