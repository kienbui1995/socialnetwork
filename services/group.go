package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateGroup func
func CreateGroup(userid int, group models.Group) (int, error) {

	p := neoism.Props{
		"name":        group.Name,
		"description": group.Description,
		"avatar":      group.Avatar,
		"cover":       group.Cover,
		"privacy":     group.Privacy,
		"status":      group.Status,
		"members":     0,
		"posts":       0,
	}
	stmt := `
	    MATCH(u:User) WHERE ID(u) = {fromid}
	  	CREATE (s:Group { props } )<-[r:CREATE]-(u)
			SET s.created_at = TIMESTAMP()
			RETURN ID(s) as id
	  	`

	params := map[string]interface{}{
		"props":  p,
		"fromid": userid,
	}
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

// GetGroups func
func GetGroups(userid int, orderby string, skip int, limit int, typegroup int) ([]models.GroupJoin, error) {

	stmt := fmt.Sprintf(`
    MATCH(u:User) WHERE ID(u) = {userid}
  	MATCH (g:Group)
		RETURN
			ID(s) AS id, s.name AS name, s.description AS description, s.avatar AS avatar, s.cover AS cover,
			s.members AS members, s.posts AS posts,
			s.created_at AS created_at, s.updated_at AS updated_at,
			CASE s.privacy when null then 1 else s.privacy end AS privacy,
			CASE s.status when null then 1 else s.status end AS status,
			exists(u-[:JOIN]->(g)) AS is_joined
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
  	`, orderby)

	params := map[string]interface{}{
		"userid": userid,
		"skip":   skip,
		"limit":  limit,
	}
	res := []models.GroupJoin{}
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

// GetGroup func
func GetGroup(groupid int, myuserid int) (models.GroupJoin, error) {
	return models.GroupJoin{}, nil
}

// UpdateGroup func
func UpdateGroup(group models.Group) (models.Group, error) {
	var stmt string
	var params map[string]interface{}
	stmt = `
			MATCH (s:Group)<-[r:JOIN]-(u:User)
			WHERE ID(s) = {groupid}
			SET s.name = {name}, s.description = {description}, s.avatar = {avatar}, s.cover = {cover}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}
			RETURN
				ID(s) AS id, s.name AS name, s.description AS description,
				s.avatar AS avatar, s.cover AS cover,
				s.posts AS posts, s.members AS members,
				s.created_at AS created_at, s.updated_at AS updated_at,
				case s.privacy when null then 1 else s.privacy end AS privacy,
				case s.status when null then 1 else s.status end AS status
			`
	params = map[string]interface{}{
		"groupid":     group.ID,
		"name":        group.Name,
		"description": group.Description,
		"avatar":      group.Avatar,
		"cover":       group.Cover,
		"privacy":     group.Privacy,
		"status":      group.Status,
	}

	res := []models.Group{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Group{}, err
	}
	if len(res) > 0 && res[0].ID >= 0 {
		return res[0], nil
	}

	return models.Group{}, errors.New("Dont' update group information")
}

// DeleteGroup func
func DeleteGroup(groupid int) (bool, error) {
	stmt := `
  	MATCH (s:Group) WHERE ID(s) = {groupid}
		DETACH DELETE s
  	`
	params := map[string]interface{}{"groupid": groupid}
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

// CheckExistGroup func
func CheckExistGroup(groupid int) (bool, error) {
	stmt := `
	MATCH (g:Group) WHERE ID(g)={groupid} RETURN ID(g) AS id;
	`
	params := neoism.Props{"groupid": groupid}

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

	if len(res) > 0 && res[0].ID == groupid {
		return true, nil
	}
	return false, nil
}

// GetGroupFeed func to get post list that group has
func GetGroupFeed(groupid int, myuserid int) ([]models.Post, error) {
	stmt := `
		MATCH(me:User) WHERE ID(me) = {myuserid}
		MATCH (g:Group)-[:HAD]->(s:Post)<-[:POST]-(u:User)
		WHERE ID(g) = {groupid}
		RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.photo when null then "" else s.photo end AS photo,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			s.likes AS likes, s.comments AS comments, s.shares AS shares,
			exists((me)-[:LIKE]->(s)) AS is_liked,
			CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_edit,
			CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_delete
		`
	params := map[string]interface{}{"groupid": groupid, "myuserid": myuserid}
	res := []models.Post{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 && res[0].PostID >= 0 {
		return res, nil

	}

	return nil, errors.New("ERROR in GetUserPost service: Groupid <0")
}

// GetJoinedGroups func
func GetJoinedGroups(userid int, orderby string, skip int, limit int) ([]models.Group, error) {

	stmt := fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
		  	MATCH (g:Group)<-[r:JOIN]-(u)
				RETURN
					ID(s) AS id, s.name AS name, s.description AS description, s.avatar AS avatar, s.cover AS cover,
					s.members AS members, s.posts AS posts,
					s.created_at AS created_at, s.updated_at AS updated_at,
					CASE s.privacy when null then 1 else s.privacy end AS privacy,
					CASE s.status when null then 1 else s.status end AS status
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, orderby)

	params := map[string]interface{}{
		"userid": userid,
		"skip":   skip,
		"limit":  limit,
	}
	res := []models.Group{}
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

// CreateJoinGroup func
func CreateJoinGroup(groupid int, userid int) (bool, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u) = {userid}
	MATCH(g:Group) WHERE ID(g) = {groupid}
	MERGE(u)-[l:JOIN]->(g)
	ON CREATE SET l.created_at = TIMESTAMP()
	RETURN exists((u)-[l]->(g)) AS joined
	`
	params := map[string]interface{}{"groupid": groupid, "userid": userid}
	res := []struct {
		Joined bool `json:"joined"`
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
	if len(res) > 0 && res[0].Joined == true {
		return true, nil
	}
	return false, nil
}

// GetGroupMembers func
func GetGroupMembers(groupid int, myuserid int) ([]models.SUser, error) {
	stmt := `
    MATCH (u:User)-[r:JOIN]->(g:Group) WHERE ID(g) = {groupid} WITH u
		MATCH (me:User) WHERE ID(me) = {myuserid}

		RETURN
			ID(u) AS id,
			u.username AS username,
			u.full_name AS full_name,
			u.avatar AS avatar,
			exists((me)-[:FOLLOW]->(u)) AS is_followed
  	`
	params := map[string]interface{}{
		"groupid":  groupid,
		"myuserid": myuserid,
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

// DeleteJoinGroup func
func DeleteJoinGroup(groupid int, userid int) (bool, error) {
	stmt := `
	MATCH (u:User)-[l:JOIN]->(g:Group) WHERE ID(g) = {groupid} AND ID(u) = {userid}
	DELETE l
	`
	params := map[string]interface{}{"groupid": groupid, "userid": userid}

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

// CheckJoinGroup func
func CheckJoinGroup(groupid int, userid int) (bool, error) {
	stmt := `
  	MATCH (u:User)-[l:JOIN]->(g:Group)
		WHERE ID(u) = {userid} AND ID(g) = {groupid}
		RETURN ID(l) as id
  	`
	params := neoism.Props{"groupid": groupid, "userid": userid}
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

// IncreaseGroupMembers func
func IncreaseGroupMembers(groupid int) (bool, error) {
	stmt := `
	MATCH (g:Group)
	WHERE ID(g)= {groupid}
	SET g.members = g.members+1
	RETURN ID(g) AS id
	`
	params := neoism.Props{"groupid": groupid}
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
	if len(res) > 0 && res[0].ID == groupid {
		return true, nil
	}
	return false, nil
}

// DecreaseGroupMembers func
func DecreaseGroupMembers(groupid int) (bool, error) {
	stmt := `
	MATCH (g:Group)
	WHERE ID(g)= {groupid}
	SET g.members = g.members-1
	RETURN ID() AS id
	`
	params := neoism.Props{"groupid": groupid}
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
	if len(res) > 0 && res[0].ID == groupid {
		return true, nil
	}
	return false, nil
}
