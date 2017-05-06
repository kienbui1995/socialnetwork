package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUserStatus func
func CreateUserStatus(userid int, message string, privacy int, status int) (int, error) {
	if userid >= 0 && len(message) > 0 {
		p := neoism.Props{
			"message":  message,
			"privacy":  privacy,
			"status":   status,
			"likes":    0,
			"comments": 0,
			"shares":   0,
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

//GetUserStatuses func
func GetUserStatuses(userid int, myuserid int, orderby string, skip int, limit int) ([]models.UserStatus, error) {

	stmt := fmt.Sprintf(`
	    MATCH(u:User) WHERE ID(u) = {userid}
			MATCH(me:User) WHERE ID(me) = {myuserid}
	  	MATCH (s:Status)<-[r:POST]-(u)
			WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
			RETURN
				ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
				s.likes AS likes, s.comments AS comments, s.shares AS shares,
				ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username,
				exists((me)-[:LIKE]->(s)) AS is_liked
			ORDER BY %s
			SKIP {skip}
			LIMIT {limit}
	  	`, orderby)
	params := map[string]interface{}{
		"userid":   userid,
		"myuserid": myuserid,
		"skip":     skip,
		"limit":    limit,
	}
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
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			exists((u)-[:LIKE]->(s)) AS is_liked
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
	stmt := `
  	MATCH (s:Status)<-[r:POST]-(u:User) WHERE ID(s) = {statusid}
		DETACH DELETE s
  	`
	params := map[string]interface{}{"statusid": statusid}
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

// GetUserIDPostedStatus func
func GetUserIDPostedStatus(statusid int) (int, error) {
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
	return -1, nil
}

// CheckExistUserStatus func
func CheckExistUserStatus(statusid int) (bool, error) {
	return CheckExistNodeWithID(statusid)
}

// GetUserStatus func
func GetUserStatus(statusid int, myuserid int) (models.UserStatus, error) {
	if statusid >= 0 {

		stmt := `
		MATCH(me:User) WHERE ID(me) = {myuserid}
		MATCH (s:Status)<-[:POST]-(u:User)
		WHERE ID(s) = {statusid}
		RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			s.likes AS likes, s.comments AS comments, s.shares AS shares,
			exists((me)-[:LIKE]->(s)) AS is_liked
		`
		params := map[string]interface{}{"statusid": statusid, "myuserid": myuserid}
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
func CreateStatusLike(statusid int, userid int) (int, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u) = {userid}
	MATCH(s:Status) WHERE ID(s) = {statusid}
	MERGE(u)-[l:LIKE]->(s)
	ON CREATE SET l.created_at = TIMESTAMP()
	RETURN exists((u)-[l]->(s)) AS liked, s.likes AS likes
	`
	params := map[string]interface{}{"statusid": statusid, "userid": userid}
	res := []struct {
		Liked bool `json:"liked"`
		Likes int  `json:"likes"`
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
	if len(res) > 0 && res[0].Liked == true {
		return res[0].Likes + 1, nil
	}
	return -1, nil
}

// GetStatusLikes func
func GetStatusLikes(statusid int, myuserid int, orderby string, skip int, limit int) ([]models.SUserLike, error) {

	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {myuserid}
	MATCH (u:User)-[l:LIKE]->(s:Status)
	WHERE ID(s) = {statusid}
	RETURN
		ID(u) AS id, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
		l.created_at as liked_at,
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
	res := []models.SUserLike{}
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
func DeleteStatusLike(statusid int, userid int) (int, error) {
	stmt := `
	MATCH (u:User)-[l:LIKE]->(s:Status) WHERE ID(s) = {statusid} AND ID(u) = {userid}
	DELETE l
	RETURN s.likes AS likes
	`
	params := map[string]interface{}{"statusid": statusid, "userid": userid}
	res := []struct {
		Likes int `json:"likes"`
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
	if len(res) > 0 {
		return res[0].Likes - 1, nil
	}
	return -1, nil
}

// CheckExistStatusLike func
func CheckExistStatusLike(statusid int, userid int) (bool, error) {
	stmt := `
  	MATCH (u:User)-[l:LIKE]->(s:Status)
		WHERE ID(u) = {userid} AND ID(s) = {statusid}
		RETURN ID(l) as id
  	`
	params := neoism.Props{"statusid": statusid, "userid": userid}
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
	if len(res) > 0 && res[0].ID >= 0 {
		return true, nil
	}
	return false, nil
}

// IncreaseStatusLikes func
func IncreaseStatusLikes(statusid int) (bool, error) {
	stmt := `
	MATCH (s:Status)
	WHERE ID(s)= {statusid}
	SET s.likes = s.likes+1
	RETURN ID(s) AS id
	`
	params := neoism.Props{"statusid": statusid}
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
	return false, nil
}

// DecreaseStatusLikes func
func DecreaseStatusLikes(statusid int) (bool, error) {
	stmt := `
	MATCH (s:Status)
	WHERE ID(s)= {statusid}
	SET s.likes = s.likes-1
	RETURN ID(s) AS id
	`
	params := neoism.Props{"statusid": statusid}
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
	return false, nil
}

// CheckStatusInteractivePermission func to check interactive permisson for user with a status
func CheckStatusInteractivePermission(statusid int, userid int) (bool, error) {
	stmt := `
		MATCH (who:User) WHERE ID(who) = {userid}
		MATCH (u:User)-[r:POST]->(s:Status)
		WHERE ID(s) = {statusid}
		RETURN exists((who)-[:FOLLOW]->(u)) AS followed, s.privacy AS privacy
		`
	params := map[string]interface{}{"userid": userid, "statusid": statusid}
	res := []struct {
		Followed bool `json:"followed"`
		Privacy  int  `json:"privacy"`
		Owner    bool `json:"owner"`
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
	if len(res) > 0 {
		if res[0].Privacy == configs.Public || (res[0].Followed && res[0].Privacy == configs.ShareToFollowers || res[0].Owner) {
			return true, nil
		}
	}
	return false, nil

}

//
// func GetUserInterestedStatus(statusid int) ([]models.SUser, error) {
// 	stmt := `
// 		MATCH(u:Status) WHERE ID(s) = {statusid}
// 		MATCH (u:User)-[:LIKE]-->(s)
// 		WHERE  s.privacy = 1 OR (s.privacy = 2 AND exists((who)-[:FOLLOW]->(u))) OR who = u
// 		RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
// 		exists((a)-[:FOLLOW]->(u)) as is_followed
// 		`
// 	params := map[string]interface{}{"statusid": statusid}
// 	res := []struct {
// 		Followed bool `json:"followed"`
// 		Privacy  int  `json:"privacy"`
// 		Owner    bool `json:"owner"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return false, err
// 	}
// 	if len(res) > 0 {
// 		if res[0].Privacy == configs.Public || (res[0].Followed && res[0].Privacy == configs.ShareToFollowers || res[0].Owner) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
//
// }
