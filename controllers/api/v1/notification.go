package v1

import (
	"fmt"
	"log"

	"github.com/maddevsio/fcm"
)

//PushNotification func
func PushNotification(deviceid []string) (bool, error) {
	data := map[string]string{
		"id":  "noti",
		"sum": "Happy Day",
	}
	c := fcm.NewFCM("AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V")
	token := "d2dvP8sjYUI:APA91bEyxraiHo-UKMeBAx-Pt7Mveih2Ydd1dddRK8lxbw-3_gZ78kz3uJWRdRTVgzlKp5_yumpn7dIjIkVoEBWbBRJZaHDJfYt2ydp0atLgfHcyQkOAuNCdEvK_uMM1bBZA7ayyx6HM"
	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  []string{token},
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
