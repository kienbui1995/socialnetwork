package v1

import (
	"fmt"
	"strconv"

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

		return
	}
	id, err := services.Login(json)

	if err != nil {
		if id >= 0 {
			libs.ResponseAuthJSON(c, 412, "Error login: "+err.Error())

		} else {
			libs.ResponseAuthJSON(c, 409, "No exist user: "+err.Error())
		}
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
	data := map[string]interface{}{"token": tokenstring, "id": id}
	libs.ResponseSuccessJSON(c, 1, "Login successful!", data)
}

// LoginViaFacebook func is login or sign up via Facebook
func LoginViaFacebook(c *gin.Context) {
	type FacebookToken struct {
		ID          string `json:"id"`
		AccessToken string `json:"access_token"`
		Device      string `json:"device"`
	}
	var json FacebookToken
	err := c.Bind(&json)
	if err != nil {
		c.Abort()
		return
	}
	if len(json.ID) == 0 || len(json.Device) == 0 || len(json.AccessToken) == 0 {
		libs.ResponseAuthJSON(c, 101, "Missing a few fields.")
		return
	}

	if id, errExist := services.CheckExistFacebookID(json.ID); errExist == nil && id != 0 {
		verify := libs.VerifyFacebookID(json.ID, json.AccessToken)
		if verify == true {
			tokenstring, errtoken := middlewares.GenerateToken(id, json.Device, SuperSecretPassword)
			if errtoken != nil {
				libs.ResponseAuthJSON(c, 408, "Error in generate token: "+errtoken.Error())
				return
			}
			if saveToken, err := services.SaveToken(id, json.Device, tokenstring); saveToken != true || err != nil {
				libs.ResponseBadRequestJSON(c, 1, "Don't save token"+err.Error())
				return
			}
			data := map[string]interface{}{"token": tokenstring, "id": id}
			libs.ResponseSuccessJSON(c, 1, "Login successful!", data)
			return
		}
		libs.ResponseBadRequestJSON(c, 411, "Error in checking facebook access token.")

	} else {

		libs.ResponseNotFoundJSON(c, 410, "No exist account with this facebook.")
		return
	}
	// libs.ResponseBadRequestJSON(c, -1, "Login Facebook fail")
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
	libs.ResponseSuccessJSON(c, 1, "Logout successful", nil)
}

//ForgotPassword func
func ForgotPassword(c *gin.Context) {
	var json struct {
		Email string `json:"email"`
	}
	err := c.Bind(&json)
	if err != nil {
		c.Abort()
		return
	}

	if len(json.Email) == 0 {
		libs.ResponseAuthJSON(c, 101, "Missing a few fields.")
		c.Abort()
		return
	}
	existemail, err := services.CheckExistEmail(json.Email)
	if err != nil {
		c.Abort()
		return
	}

	if existemail == true { // send password via mail
		type RecoverPassword struct {
			Email        string `json:"email"`
			RecoveryCode string `json:"recovery_code"`
		}
		recoverpass := RecoverPassword{Email: json.Email, RecoveryCode: libs.RandNumberBytes(6)}
		if err := services.CreateRecoverPassword(recoverpass.Email, recoverpass.RecoveryCode); err != nil {
			libs.ResponseServerErrorJSON(c)
			panic("Error in creating recover password: " + err.Error())
		}
		sender := libs.NewSender("kien.laohac@gmail.com", "ytanyybkizzygqjk")
		var email []string
		email = append(email, recoverpass.Email)
		go sender.SendMail(email, fmt.Sprintf("Recover password on TLSEN"), fmt.Sprintf("\ncode: %s\n Please verify within 30 minutes.", recoverpass.RecoveryCode))
		libs.ResponseSuccessJSON(c, 1, "A email sent.", nil)
	} else { // no exist email
		libs.ResponseAuthJSON(c, 413, "No exist email.")
	}
}

//VerifyRecoveryCode func
func VerifyRecoveryCode(c *gin.Context) {
	var json struct {
		Email        string `json:"email"`
		RecoveryCode string `json:"recovery_code"`
	}
	err := c.Bind(&json)
	if err != nil {
		c.Abort()
		return
	}

	if len(json.Email) == 0 || len(json.RecoveryCode) == 0 {
		libs.ResponseAuthJSON(c, 101, "Missing a few fields.")
		c.Abort()
		return
	}
	// check exist user with this email
	exist, err := services.CheckExistEmail(json.Email)
	if err != nil {
		libs.ResponseServerErrorJSON(c)
		panic("Error in check exist email: " + err.Error())

	} else if exist == true {
		id, err := services.VerifyRecoveryCode(json.Email, json.RecoveryCode)
		if err != nil {
			libs.ResponseServerErrorJSON(c)
			panic("Error in verify recovery code: " + err.Error())

		}
		if id >= 0 { // generate a key
			key := libs.RandStringBytes(6)

			libs.ResponseSuccessJSON(c, 1, "ID user and key to create new password", map[string]interface{}{"id": id, "recovery_key": key})
			go func() {
				err := services.AddUserRecoveryKey(id, key)
				if err != nil {
					libs.ResponseServerErrorJSON(c)
					panic(err)
				}
			}()

		} else {
			libs.ResponseNotFoundJSON(c, 414, "Error in recover password: Wrong recovery code.")
		}
	} else { // user with this emmail don't exist
		libs.ResponseNotFoundJSON(c, 409, "No exist user.")
	}

}

//RenewPassword func
func RenewPassword(c *gin.Context) {
	var json struct {
		ID          string `json:"id"`
		RecoveryKey string `json:"recovery_key"`
		NewPassword string `json:"new_password"`
	}
	err := c.Bind(&json)
	if err != nil {
		libs.ResponseBadRequestJSON(c, 100, "Invalid parameter")
		return
	}

	if len(json.ID) == 0 || len(json.RecoveryKey) == 0 || len(json.NewPassword) == 0 {
		libs.ResponseAuthJSON(c, 101, "Missing a few fields.")
		c.Abort()
		return
	}
	id, err := strconv.Atoi(json.ID)
	if err != nil {
		libs.ResponseBadRequestJSON(c, 110, "Invalid user id")
		return
	}
	if ok, err := services.RenewPassword(id, json.RecoveryKey, json.NewPassword); err == nil && ok == true {
		libs.ResponseSuccessJSON(c, 1, "Renew password successful", nil)
		go func() {
			okdel, errdel := services.DeleteRecoveryProperty(id)
			if errdel != nil {
				fmt.Printf("Error in DeleteRecoveryProperty: " + errdel.Error())
			}
			if okdel != true {
				fmt.Printf("Error in DeleteRecoveryProperty: No delete")
			}
		}()
	} else if err != nil {
		libs.ResponseServerErrorJSON(c)
		panic("Error in creating new password: " + err.Error())
	} else if ok != true {
		libs.ResponseNotFoundJSON(c, 310, "User data edit failure")
	}
}
