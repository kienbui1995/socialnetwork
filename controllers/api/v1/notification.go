package v1

import (
	"fmt"

	"github.com/maddevsio/fcm"
)

//PushNotification func
func PushNotification(deviceid []string) (bool, error) {
	data := map[string]string{
		"id":  "noti",
		"sum": "Happy Day",
	}
	c := fcm.NewFCM("AAAAuET9LvY:APA91bEYl-fIkcY0w7b6umgBHD4yrZnG_v9I2iY1K3EnjUfSrYvlFYIG5vrmP8wFCH8ZMZ-Kx6U6u3XIsw-AIGehs-msWXtlzOq8R_50qAiqcsrJv9WQluALvjWPqSIAPrVS2RKZ4H6V")

	response, err := c.Send(fcm.Message{
		Data:             data,
		RegistrationIDs:  deviceid,
		ContentAvailable: true,
		Priority:         fcm.PriorityHigh,
		Notification: fcm.Notification{
			Title: "Follow",
			Body:  "World",
			Sound: "",
			Badge: "113",
		},
	})
	if err != nil {

		fmt.Println("Status Code   :", response.StatusCode)
		fmt.Println("Success       :", response.Success)
		fmt.Println("Fail          :", response.Fail)
		fmt.Println("Canonical_ids :", response.CanonicalIDs)
		fmt.Println("Topic MsgId   :", response.MsgID)
		return false, err
	}
	return true, nil
}
