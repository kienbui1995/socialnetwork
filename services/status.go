package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUserStatus func
func CreateUserStatus(userid int, message string, privacy int, status int) (int, error) {
	if userid >= 0 && len(message) > 0 {
		p := neoism.Props{
			"message": message,
			"privacy": privacy,
			"status":  status,
		}
		stmt := `
    MATCH(u:User) WHERE ID(u) = {fromid}
  	CREATE (s:Status:Post { props } )<-[r:POST]-(u) SET s.created_at = TIMESTAMP() RETURN ID(s) as id
  	`
		params := map[string]interface{}{"props": p, "fromid": userid}
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
func GetUserStatuses(userid int, orderby string, skip int, limit int) ([]models.UserStatus, error) {

	stmt := fmt.Sprintf(`
	    MATCH(u:User) WHERE ID(u) = {userid}
	  	MATCH (s:Status)<-[r:POST]-(u)
			RETURN
				ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
				ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username
			ORDER BY %s
			SKIP {skip}
			LIMIT {limit}
	  	`, orderby)
	params := map[string]interface{}{
		"userid": userid,
		"skip":   skip,
		"limit":  limit,
	}
	res := []models.UserStatus{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	fmt.Print(cq)
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 && res[0].ID >= 0 {

		return res, nil

	}
	return nil, nil
}

// UpdateUserStatus func
func UpdateUserStatus(statusid int, message string, privacy int, status int) (models.UserStatus, error) {
	if statusid >= 0 {
		stmt := `
  	MATCH (s:Status)<-[r:POST]-(u:User)
    WHERE ID(s) = {statusid}
		SET s.message = {message}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}
    RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar
  	`
		params := map[string]interface{}{"statusid": statusid, "message": message, "privacy": privacy, "status": status}
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
	return CheckExistNodeWithID(statusid)
}

// GetUserStatus func
func GetUserStatus(statusid int) (models.UserStatus, error) {
	if statusid >= 0 {

		stmt := `
		MATCH (s:Status)<-[:POST]-(u:User)
		WHERE ID(s) = {statusid}
		RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar
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
			return models.UserStatus{}, err
		}
		if len(res) > 0 && res[0].ID == statusid {
			return res[0], nil
		}
	}

	return models.UserStatus{}, errors.New("ERROR in GetUserStatus service: statusid <0")
}

// CreateStatusLike func
func CreateStatusLike(statusid int, userid int) (bool, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u) = {userid}
	MATCH(s:Status) WHERE ID(s) = {statusid}
	CREATE UNIQUE (u)-[l:LIKE{created_at:TIMESTAMP()}]->(s)
	RETURN exists(l) AS liked
	`
	params := map[string]interface{}{"statusid": statusid, "userid": userid}
	res := []struct {
		Liked bool `json:"liked"`
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
	if len(res) > 0 && res[0].Liked == true {
		return true, nil
	}
	return false, nil
}

// GetStatusLikes func
func GetStatusLikes(statusid int, myuserid int, orderby string, skip int, limit int) ([]models.SUser, error) {

	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {myuserid}
	MATCH (u:User)-[l:LIKE]->(s:Status) WHERE ID(s) = {statusid}
	RETURN ID(u) AS id, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
	exists((me)-[:FOLLOW]->(u)) AS is_followed
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, orderby)
	params := map[string]interface{}{
		"statusid": statusid,
		"myuserid": myuserid,
		"skip":     skip,
		"limit":    limit,
	}
	res := []models.SUser{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 && res[0].UserID >= 0 {
		return res, nil
	}
	return nil, nil
}

// DeleteStatusLike func
func DeleteStatusLike(statusid int, userid int) (bool, error) {
	stmt := `
	MATCH (u:User)-[l:LIKE]->(s:Status) WHERE ID(s) = {statusid} AND ID(u) = {userid}
	DELETE l
	`
	params := map[string]interface{}{"statusid": statusid, "userid": userid}
	res := []models.SUser{}
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
