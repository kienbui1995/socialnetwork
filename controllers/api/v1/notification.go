package v1

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/libs"
	"github.com/maddevsio/fcm"
)

//PushNotification func
func PushNotification(deviceid []string) (bool, error) {
	data := map[string]string{
		"id":  "noti",
		"sum": "Happy Day",
	}
	c := fcm.NewFCM("AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V")
	//token := "d2dvP8sjYUI:APA91bEyxraiHo-UKMeBAx-Pt7Mveih2Ydd1dddRK8lxbw-3_gZ78kz3uJWRdRTVgzlKp5_yumpn7dIjIkVoEBWbBRJZaHDJfYt2ydp0atLgfHcyQkOAuNCdEvK_uMM1bBZA7ayyx6HM"
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  deviceid,
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: "Follow",
			Body:  "World",
			Sound: "default",
			Badge: "113",
		},
	})
	if err != nil {
		log.Fatal(err)
		fmt.Println("Status Code   :", response.StatusCode)
		fmt.Println("Success       :", response.Success)
		fmt.Println("Fail          :", response.Fail)
		fmt.Println("Canonical_ids :", response.CanonicalIDs)
		fmt.Println("Topic MsgId   :", response.MsgID)
		return false, err
	}
	return true, nil
}

// PushTest func
func PushTest(c *gin.Context) {
	json := struct {
		TokenClient string `json:"token_client"`
		Message     string `json:"message"`
		Title       string `json:"title"`
	}{}
	if errBind := c.Bind(&json); errBind != nil {
		libs.ResponseBadRequestJSON(c, 101, "Bug push")
		return
	}
	push := fcm.NewFCM(configs.FCMToken)
	data := map[string]string{
		"id":  "noti",
		"sum": "Happy Day",
	}

	response, err := push.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{json.TokenClient},
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: json.Title,
			Body:  json.Message,
			Sound: "default",
			Badge: "113",
		},
	})
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Status Code   :", response.StatusCode)
		fmt.Println("Success       :", response.Success)
		fmt.Println("Fail          :", response.Fail)
		fmt.Println("Canonical_ids :", response.CanonicalIDs)
		fmt.Println("Topic MsgId   :", response.MsgID)
		libs.ResponseBadRequestJSON(c, -1, err.Error())
		return
	}
	libs.ResponseJSON(c, 200, 1, "push: "+json.Message+json.Title+json.TokenClient, nil)

}
