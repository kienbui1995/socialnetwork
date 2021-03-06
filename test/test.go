package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

//const
const (
	GraphFacebookAPI = "https://graph.facebook.com"
)

//VerifyFacebookID func to check logged in via Facebook
func VerifyFacebookID(id string, accessToken string) string {
	url := fmt.Sprintf("%s/me?fields=id&access_token=%s", GraphFacebookAPI, accessToken)
	res, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	// read json http response
	jsonDataFromHTTP, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}

	var jsonData struct{ id string }

	err = json.Unmarshal([]byte(jsonDataFromHTTP), &jsonData) // here!

	if err != nil {
		panic(err)
	}

	return jsonData.id

}

func main() {
	router := gin.Default()
	test := router.Group("/")
	{
		test.POST("vertify_facebook_id", handlers)
	}

}
