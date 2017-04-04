package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUser func
func CreateUser(user models.User) (models.User, error) {
	stmt := `
	Create (u:User{
		Username: {username},
		Password: {password},
		Firstname: {firstname},
		Lastname: {lastname},
		Fullname: {fullname},
		Gender: {gender},
		Birthday: {birthday},
		FacebookId: {facebookid},
		FacebookToken: {facebooktoken},
		Email: {email},
		Avatar: {avatar},
		Status: {status}
		})
	return u as user
	`
	params := neoism.Props{
		"username":      user.Data["username"],
		"password":      user.Data["password"],
		"firstname":     user.Data["firstname"],
		"lastname":      user.Data["lastname"],
		"fullname":      user.Data["fullname"],
		"birthday":      user.Data["birthday"],
		"gender":        user.Data["gender"],
		"facebookid":    user.Data["facebookid"],
		"facebooktoken": user.Data["facebooktoken"],
		"email":         user.Data["email"],
		"avatar":        user.Data["avatar"],
		"status":        user.Data["status"],
	}
	res := []map[string]map[string]map[string]interface{}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err == nil {
		var userreturn models.User
		userreturn.Data = make(map[string]interface{})
		for k, v := range res[0]["user"] {
			if v != nil {
				if k == "data" || k == "metadata" {
					for k1, v1 := range v {
						if k1 != "labels" {
							if k1 == "id" {
								k1 = "userid"
							}
							userreturn.Data[strings.ToLower(k1)] = v1
						}
					}

				}
			}
		}

		return userreturn, nil
	}
	return user, err
}

// GetAllUser func
func GetAllUser() ([]map[string]interface{}, error) {
	user := models.User{}
	stmt := `
	MATCH (u:User) RETURN ID(u) as userid, u.Username as username, u.Password as password, u.Email as email, u.Status as status, u.Avatar as avatar LIMIT 25;
	`
	res := []map[string]interface{}{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	k := int(len(res))

	if k < 1 {
		return nil, errors.New("No User")
	}
	var listuser []map[string]interface{}

	for i := range res {
		user.Data = make(map[string]interface{})
		for k, v := range res[i] {
			if v != nil {
				user.Data[k] = v
			}
		}
		listuser = append(listuser, user.Data)
	}
	return listuser, nil
}

// GetUser func
func GetUser(userid int) (models.User, error) {
	var user = models.User{}
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userId} RETURN u.Avatar as avatar, ID(u) as userid, u.Username as username, u.Password as password, u.Email as email, u.Status as status LIMIT 25;
	`
	params := neoism.Props{"userId": userid}

	res := []map[string]interface{}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)

	if err != nil {
		return user, err
	}
	if len(res) == 1 {
		user.Data = make(map[string]interface{})
		for k, v := range res[0] {
			if v != nil {
				user.Data[k] = v
			}
		}

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
	MATCH (u:User) WHERE ID(u) = {userid} SET u.Email = {email}, u.Status = {status},  u.Avatar = {avatar}, u.Password = {password};
	`
	params := neoism.Props{"userid": user.Data["userid"], "email": user.Data["email"], "status": user.Data["status"], "avatar": user.Data["avatar"], "password": user.Data["password"]}
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
