package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// Login func to user login system
func Login(login models.Login) (int, error) {
	stmt := `
	MATCH (u:User) WHERE u.username	 = {username} return ID(u) as id, u.password as password
	`
	params := neoism.Props{"username": login.Username, "password": login.Password}

	res := []struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
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
	} else if res[0].Password == login.Password {
		return res[0].ID, nil
	}
	return res[0].ID, errors.New("Wrong password")
}

// SaveToken func t insert token to db
func SaveToken(userid int, device string, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
		MERGE (u)-[:LOGGED_IN]->(d:Device {device:{device}}) SET d.token = {token}
	` // chua test
	params := neoism.Props{"userid": userid, "token": token, "device": device}

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
	 MATCH (u:User) WHERE ID(u) = {userid} return exists( ((u)-[:LOGGED_IN]->(:Device{ token:{token}})) ) as existToken
	`
	params := neoism.Props{"userid": userid, "token": token}

	res := []struct {
		ExistToken bool `json:"existToken"`
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
		return false, errors.New("Error, can't check token on DB")
	}
	if res[0].ExistToken != true {
		return false, errors.New("Wrong token")
	}
	return true, nil
}

//DeleteToken func to delete token of user
func DeleteToken(userid int, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} MATCH ((u)-[:LOGGED_IN]->(d)) WHERE d.token = {token} DETACH DELETE d
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
