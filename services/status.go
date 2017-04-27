package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUserStatus func
func CreateUserStatus(status models.UserStatus) (int, error) {
	if status.UserID >= 0 && len(status.Message) > 0 {
		p := neoism.Props{
			"message": status.Message,
		}
		stmt := `
    MATCH(u:User) WHERE ID(u) = {fromid}
  	CREATE (s:Status:Post { props } )<-[r:POST]-(u) SET s.created_at = TIMESTAMP() RETURN ID(s) as id
  	`
		params := map[string]interface{}{"props": p, "fromid": status.UserID}
		res := []struct {
			ID int `json:"id"`
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
		if len(res) > 0 && res[0].ID >= 0 {
			return res[0].ID, nil
		}

	}
	return -1, nil
}

// GetUserStatuses func
func GetUserStatuses(userid int) ([]models.UserStatus, error) {
	if userid >= 0 {
		stmt := `
    MATCH(u:User) WHERE ID(u) = {userid}
  	MATCH (s:Status)<-[r:POST]-(u) RETURN ID(s) AS id, s.message AS message, s.created_at AS created_at, ID(u) AS userid
  	`
		params := map[string]interface{}{"userid": userid}
		res := []models.UserStatus{}
		cq := neoism.CypherQuery{
			Statement:  stmt,
			Parameters: params,
			Result:     &res,
		}
		err := conn.Cypher(&cq)
		if err != nil {
			return nil, err
		}
		if len(res) > 0 && res[0].ID >= 0 {
			return res, nil
		}

	}
	return nil, nil
}

// UpdateUserStatus func
func UpdateUserStatus(statusid int, message string) (models.UserStatus, error) {
	if statusid >= 0 {
		stmt := `
  	MATCH (s:Status)<-[r:POST]-(u:User)
    WHERE ID(s) = {statusid}  SET s.message = {message}, s.updated_at = TIMESTAMP()
    RETURN ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at, ID(u) AS userid
  	`
		params := map[string]interface{}{"statusid": statusid, "message": message}
		res := []models.UserStatus{}
		cq := neoism.CypherQuery{
			Statement:  stmt,
			Parameters: params,
			Result:     &res,
		}
		err := conn.Cypher(&cq)
		if err != nil {
			return models.UserStatus{}, err
		}
		if len(res) > 0 && res[0].ID >= 0 {
			return res[0], nil
		}

	}
	return models.UserStatus{}, errors.New("Dont' update user status")
}

// DeleteUserStatus func
func DeleteUserStatus(statusid int) (bool, error) {
	if statusid >= 0 {
		stmt := `
  	MATCH (s:Status)<-[r:POST]-(u:User) WHERE ID(s) = {statusid} DELETE r, s
  	`
		params := map[string]interface{}{"statusid": statusid}
		res := []models.UserStatus{}
		cq := neoism.CypherQuery{
			Statement:  stmt,
			Parameters: params,
			Result:     &res,
		}
		err := conn.Cypher(&cq)
		if err != nil {
			return false, err
		}

	}
	return true, nil
}

// GetUserIDPostedStatus func
func GetUserIDPostedStatus(statusid int) (int, error) {
	if statusid >= 0 {

		stmt := `
    MATCH (u:User)-[r:POST]->(s:Status) WHERE ID(s) = {statusid} RETURN ID(u) AS id
  	`
		params := map[string]interface{}{"statusid": statusid}
		res := []struct {
			ID int `json:"id"`
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
		if len(res) > 0 && res[0].ID >= 0 {
			return res[0].ID, nil
		}

	}
	return -1, nil
}

// CheckExistUserStatus func
func CheckExistUserStatus(statusid int) (bool, error) {
	if statusid >= 0 {
		stmt := `
  	MATCH (s:Status) WHERE ID(s) = {statusid} RETURN ID(s) as id
  	`
		params := map[string]interface{}{"statusid": statusid}
		res := []struct {
			ID int `json:"id"`
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
		if len(res) > 0 && res[0].ID == statusid {
			return true, nil
		}
	}
	return false, nil
}
