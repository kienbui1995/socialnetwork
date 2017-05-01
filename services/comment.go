package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateStatusComment func
func CreateStatusComment(statusid int, userid int, message string) (int, error) {
	p := neoism.Props{
		"message": message,
		"status":  1,
	}
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
	MATCH (s:Status) WHERE ID(s) = {statusid}
	CREATE (c:Comment { props } ) SET c.created_at = TIMESTAMP()
	CREATE (u)-[w:WRITE]->(c)-[a:AT]->(s)
	RETURN ID(c) AS id
	`
	params := map[string]interface{}{"props": p, "statusid": statusid, "userid": userid}
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

// GetStatusComments func
func GetStatusComments(statusid int, orderby string, skip int, limit int) ([]models.UserComment, error) {

	stmt := fmt.Sprintf(`
	MATCH (u:User)-[w:WRITE]->(c:Comment)-[a:AT]->(s:Status)
	WHERE ID(s) = {statusid}
	RETURN
		ID(c) AS id, c.message AS message, c.created_at AS created_at, c.updated_at AS updated_at ,c.status AS status,
		ID(u) AS userid, u.username AS username. u.full_name AS full_name, u.avatar AS avatar
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, orderby)
	params := map[string]interface{}{
		"statusid": statusid,
		"skip":     skip,
		"limit":    limit,
	}

	res := []models.UserComment{}
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

// DeleteStatusComment func
func DeleteStatusComment(commentid int) (bool, error) {
	stmt := `
	MATCH (c:Comment) WHERE ID(c) = {commentid}
	DETACH DELETE c
	`
	params := map[string]interface{}{"commentid": commentid}
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

// UpdateStatusComment func
func UpdateStatusComment(commentid int, message string) (bool, error) {
	stmt := `
	MATCH (c:Comment) WHERE ID(c) = {commentid}
	SET c.message = {message}, c.updated_at = TIMESTAMP()
  RETURN ID(c) AS id
	`
	params := map[string]interface{}{"commentid": commentid, "message": message}
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
	if len(res) > 0 && res[0].ID == commentid {
		return true, nil
	}
	return false, nil
}

// CheckExistComment func
func CheckExistComment(commentid int) (bool, error) {
	return CheckExistNodeWithID(commentid)
}

// GetUserIDWroteComment func
func GetUserIDWroteComment(commentid int) (int, error) {
	stmt := `
    MATCH (u:User)-[w:WRITE]->(c:Comment) WHERE ID(c) = {commentid} RETURN ID(u) AS id
  	`
	params := map[string]interface{}{"commentid": commentid}
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
