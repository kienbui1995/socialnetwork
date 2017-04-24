package v1

import (
	"errors"

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

// GetUserIDFromToken func return userid to check permission
func GetUserIDFromToken(token string) (int, error) {

	if len(token) == 0 {
		return -1, errors.New("NULL userid in token")
	}

	claims, errclaim := middlewares.ExtractClaims(token, SuperSecretPassword)
	if errclaim != nil {
		return -1, errclaim
	}
	return int(claims["userid"].(float64)), nil
}

// CheckPermissionsWithID func return true if id = userid in token
func CheckPermissionsWithID(id int, token string) (bool, error) {
	if len(token) == 0 || id <= 0 {
		return false, errors.New("ID or Token is NULL in checking permission")
	}

	claims, errclaim := middlewares.ExtractClaims(token, SuperSecretPassword)
	if errclaim != nil {
		return false, errclaim
	}
	if id == int(claims["userid"].(float64)) {
		return true, nil
	}
	return false, nil
}
