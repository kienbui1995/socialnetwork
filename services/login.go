package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// Login func to user login system
func Login(login models.Login) (int, error) {
	stmt := `
	MATCH (u:User) WHERE u.Username = {username} and u.Password = {password} return u.userId
	`
	params := neoism.Props{"username": login.Username, "password": login.Password}

	res := []struct {
		UserID int `json:"u.userId"`
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
		return -1, errors.New("No exist user!")
	}
	return res[0].UserID, nil
}

// SaveToken func t insert token to db
func SaveToken(userid int, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE u.userId = {userid} SET u.Token = {token}
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
