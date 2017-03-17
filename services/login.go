package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// Login func to user login system
func Login(login models.Login) (int, error) {
	stmt := `
	MATCH (u:User) WHERE u.Username	 = {username} and u.Password = {password} return ID(u)
	`
	params := neoism.Props{"username": login.Username, "password": login.Password}

	res := []struct {
		UserID int `json:"ID(u)"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err
	}
	if len(res) == 0 {
		return -1, errors.New("No exist user")
	}
	return res[0].UserID, nil
}

// SaveToken func t insert token to db
func SaveToken(userid int, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} SET u.Token = {token}
	`
	params := neoism.Props{"userid": userid, "token": token}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckExistToken func to check exist token in DB
func CheckExistToken(userid int, token string) (bool, error) {
	//check exist token
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} RETURN u.Token
	`
	params := neoism.Props{"userid": userid}

	res := []struct {
		Token string `json:"u.Token"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, errors.New("No exist token")
	}
	if res[0].Token != token {
		return false, errors.New("Wrong token")
	}
	return true, nil
}

//DeleteToken func to delete token of user
func DeleteToken(userid int) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} REMOVE u.Token
	`
	params := neoism.Props{"userid": userid}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}
