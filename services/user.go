package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUser func
func CreateUser(user models.User) (int, error) {
	p := neoism.Props{
		//"userId":   int(user.UserID),
		"Username": user.Username,
		"Password": user.Password,

		"Email":  user.Email,
		"Status": int(user.Status),
	}

	node, errNode := conn.CreateNode(p)
	if errNode != nil {
		return -1, errNode
	}
	errLabel := node.SetLabels([]string{"User"})
	if errLabel != nil {
		node.Delete()
		return -1, errLabel
	}
	// var propNode neoism.Props
	// var err error
	// propNode, err = node.Properties()
	// if err != nil {
	// 	return -1, err
	// }
	userid := node.Id()

	return userid, nil
}

// GetAllUser func
func GetAllUser() ([]models.User, error) {
	user := models.User{}
	stmt := `
	MATCH (u:User) RETURN ID(u), u.Username, u.Password, u.Email, u.Status LIMIT 25;
	`
	res := []struct {
		//id       int     `json:ID(u)`
		UserID   int    `json:"ID(u)"`
		Username string `json:"u.Username"`
		Password string `json:"u.Password"`
		Email    string `json:"u.Email"`
		Status   int    `json:"u.Status"`
	}{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	k := int(len(res))
	fmt.Printf("%d\n", k)
	if k < 1 {
		return nil, errors.New("No User")
	}
	var listuser []models.User
	// fmt.Printf("%d\n", k)
	for i := range res {

		user.UserID = int(res[i].UserID)
		user.Username = res[i].Username

		user.Password = res[i].Password
		user.Status = int(res[i].Status)

		listuser = append(listuser, user)
		fmt.Printf("%d\n", i)
	}
	return listuser, nil
}

// GetUser func
func GetUser(userid int) (models.User, error) {
	var user = models.User{}
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} RETURN ID(u), u.Username, u.Password, u.Email, u.Status LIMIT 25;
	`
	params := neoism.Props{"userid": userid}
	res := []struct {
		//id       int     `json:ID(u)`
		UserID   int    `json:"ID(u)"`
		Username string `json:"u.Username"`
		Password string `json:"u.Password"`
		Email    string `json:"u.Email"`
		Status   int    `json:"u.Status"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	// k := len(res)
	// // fmt.Printf("%d\n", k)
	// for index := 0; index < k; index++ {
	// 	fmt.Printf("id: %d,\nuname: %s,\npass: %s,\nemail: %s,\nstatus: %d\n",
	// 		int(res[index].UserID), res[index].Username, res[index].Password, res[index].Email, int(res[index].Status))
	// }
	if err != nil {
		return user, err
	}
	if len(res) == 1 {
		user = models.User{UserID: res[0].UserID, Username: res[0].Username, Password: res[0].Password, Email: res[0].Email, Status: res[0].Status}
		return user, nil
	} else if len(res) > 1 {
		return user, errors.New("Many User")
	} else {
		return user, errors.New("No User")
	}
}

// UpdateUser func
func UpdateUser(user models.User) (bool, error) {

	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} SET u.Email = {email}, u.Status = {status};
	`
	params := neoism.Props{"userid": user.UserID, "email": user.Email, "status": user.Status}
	res := false
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// DeleteUser func
func DeleteUser(userid int) (int, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} delete u RETURN count(u) ;
	`
	params := neoism.Props{"userid": userid}
	res := -1
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	return res, err
}

//CheckExistUser func to check exist User
func CheckExistUser(userid int) (bool, error) {
	where := fmt.Sprintf("ID(u) = %d", userid)
	number, err := CheckExistNode("User", where)
	if err != nil {
		return false, err
	}
	if number > 0 {
		return true, nil
	}
	return false, nil
}
