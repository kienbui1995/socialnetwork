package main

import (
	"fmt"
	"log"
	//"socialnetwork/models"
	"strconv"

	"github.com/gin-gonic/gin"
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

var (
	neo4jURL = "Bolt://neo4j:tlis2016@148.72.245.101:7687"
)

//
// func test() string {
// 	return models.testUser()
// }

// func connectBolt() {
// 	driver := bolt.NewDriver()
// 	conn, err := driver.OpenNeo(neo4jURL)
//
// }

func interfaceSliceToString(s []interface{}) []string {
	o := make([]string, len(s))
	for idx, item := range s {
		o[idx] = item.(string)
	}
	return o
}

func main() {
	router := gin.Default()
	router.GET("/user", userGetList)       // get all user
	router.GET("/user/:userid", userGet)   // get a user
	router.POST("/user", userCreate)       // create a user
	router.DELETE("/user/:id", userDelete) // delete a user
	router.PUT("/user/:id", userUpdate)    // update a user
	router.Run(":8080")
}

func userGet(c *gin.Context) {

	userid, _ := strconv.ParseInt(c.Param("userid"), 10, 32)
	//fmt.Printf(" lay para user%T: %d\n", userid, userid)

	//connectBolt
	//driver := bolt.NewDriver()
	conn, err := bolt.NewDriver().OpenNeo(neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err,
		})
		//c.Write([]byte("An error occurred connecting to the DB"))
		return
	}
	defer conn.Close()
	// stmt, err := conn.PrepareNeo("MATCH (u:User) WHERE u.Username = {userid} RETURN u")
	// if err != nil {
	// 	panic(err)
	// }

	data, _, _, err := conn.QueryNeoAll("MATCH (u:User) WHERE ID(u) = {userid} RETURN u", map[string]interface{}{"userid": userid})
	if err != nil {
		panic(err)
	}
	// /fmt.Printf(" lay para user%T lần 2: %s\n", userid, userid)
	c.JSON(200, gin.H{
		"status": "success",
		//"result": result,
		"data": data,
	})
}

func userCreate(c *gin.Context) {
	//roomid := c.Param("roomid")
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	// fmt.Printf(" lay para user: %s\n", username)
	// fmt.Printf(" lay para pass: %s\n", password)
	// fmt.Printf(" lay para email: %s\n", email)
	//connectBolt
	conn, err := bolt.NewDriver().OpenNeo(neo4jURL)
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		c.JSON(500, gin.H{
			"status":  "fail",
			"message": err,
		})
		//c.Write([]byte("An error occurred connecting to the DB"))
		return
	}
	defer conn.Close()
	stmt, err := conn.PrepareNeo("CREATE (user:User {Username: {username}, Password: {password}, Email: {email}, Status: 0})")
	if err != nil {
		panic(err)
	}

	result, err := stmt.ExecNeo(map[string]interface{}{"username": username, "password": password, "email": email})

	if err != nil {
		panic(err)
	}
	fmt.Printf(" send to database with user: %s, %s, %s", username, password, email)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": result,
	})
}

func userDelete(c *gin.Context) {
	userid := c.Param("id")

	//connectBolt
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(neo4jURL)
	if err != nil {
		panic(err)
	}
	stmt, err := conn.PrepareNeo("MATCH (u:User) WHERE u.UserId = {userid} DELETE u")
	if err != nil {
		panic(err)
	}

	result, err := stmt.ExecNeo(map[string]interface{}{"userid": userid})
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
		"status":  "success",
		"message": result,
	})

}

func userUpdate(c *gin.Context) {
	userid, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	email := c.PostForm("email")
	password := c.PostForm("password")
	status := c.PostForm("status")
	//connectBolt
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(neo4jURL)
	if err != nil {
		panic(err)
	}
	// stmt, err := conn.PrepareNeo("MATCH (u:User) WHERE u.UserId = {userid} RETURN u.Status")
	// if err != nil {
	// 	panic(err)
	// }
	query := "MATCH (u:User) WHERE ID(u) = {userid} SET u.Password = {password}, u.Email = {email}, u.Status = {status} RETURN u"
	data, _, _, err := conn.QueryNeoAll(
		query,
		map[string]interface{}{"userid": userid, "password": password, "email": email, "status": status})
	if err != nil {
		panic(err)
	}
	//fmt.Printf(" aaaar: %d\n, %d\n, %d\n, %d\n", data, a, b, err)
	c.JSON(200, gin.H{
		"status":  "success",
		"message": data,
	})

}

func userGetList(c *gin.Context) {
	var status int64 = -1
	var query string
	if c.Query("status") != "" {
		status, _ = strconv.ParseInt(c.Query("status"), 10, 32)
	}
	//connectBolt
	driver := bolt.NewDriver()
	conn, err := driver.OpenNeo(neo4jURL)
	if err != nil {
		panic(err)
	}
	if status < 0 {
		query = "MATCH (n:User) RETURN n"
	} else {
		query = "MATCH (n:User) WHERE n.Status = {status} RETURN n"
	}
	data, _, _, _ := conn.QueryNeoAll(query, map[string]interface{}{"status": status})
	// var result string
	// for _, row := range data {
	//
	// 	result = result + fmt.Sprintf("FIELDS: %d %f\n", data[0][0].(int64), data[0][1].(float64)) // Prints all nodes
	// 	// /fmt.Printf("FIELDS: %d %f\n", data[0][0].(int64), data[0][1].(float64))
	// }

	c.JSON(200, gin.H{
		"status":  "success",
		"message": data,
	})

}
