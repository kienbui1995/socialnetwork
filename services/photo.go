package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUserPhoto func
func CreateUserPhoto(userid int, message string, photo string, privacy int, status int) (int, error) {
	p := neoism.Props{
		"message":  message,
		"photo":    photo,
		"privacy":  privacy,
		"status":   status,
		"likes":    0,
		"comments": 0,
		"shares":   0,
	}
	stmt := `
  MATCH(u:User) WHERE ID(u) = {userid}
  CREATE (p:Photo:Post { props } )<-[r:POST]-(u) SET p.created_at = TIMESTAMP()
  RETURN ID(p) as id
  `
	params := map[string]interface{}{"props": p, "userid": userid}
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

//GetUserPhotos func
func GetUserPhotos(userid int, myuserid int, orderby string, skip int, limit int) ([]models.Photo, error) {

	stmt := fmt.Sprintf(`
	    MATCH(u:User) WHERE ID(u) = {userid}
			MATCH(me:User) WHERE ID(me) = {myuserid}
	  	MATCH (p:Photo)<-[r:UPLOAD]-(u)
      WHERE p.privacy = 1 OR (p.privacy = 2 AND exists(me-[:FOLLOW]->(u))) OR {userid} = {myuserid}
			RETURN
				ID(p) AS id, p.message AS message, p.created_at AS created_at, p.updated_at AS updated_at,
				case p.privacy when null then 1 else p.privacy end AS privacy, case p.status when null then 1 else p.status end AS status,
				p.likes AS likes, p.comments AS comments, p.shares AS shares,
				ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username,
				exists((me)-[:LIKE]->(p)) AS is_liked
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
	res := []models.Photo{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 && res[0].PhotoID >= 0 {
		return res, nil
	}
	return nil, nil
}

// UpdateUserPhoto func
func UpdateUserPhoto(photoid int, message string, privacy int, status int) (models.Photo, error) {
	stmt := `
  	MATCH (p:Photo)<-[r:UPLOAD]-(u:User)
    WHERE ID(p) = {photoid}
		SET p.message = {message}, p.privacy = {privacy}, p.updated_at = TIMESTAMP(), p.status = {status}
    RETURN
			ID(p) AS id, p.message AS message, p.created_at AS created_at, p.updated_at AS updated_at,
			case p.privacy when null then 1 else p.privacy end AS privacy, case p.status when null then 1 else p.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
      exists((u)-[:LIKE]->(p)) AS is_liked
  	`
	params := map[string]interface{}{"photoid": photoid, "message": message, "privacy": privacy, "status": status}
	res := []models.Photo{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Photo{}, err
	}
	if len(res) > 0 && res[0].PhotoID >= 0 {
		return res[0], nil
	}
	return models.Photo{}, errors.New("Dont' update user photo")
}

// DeleteUserPhoto func
func DeleteUserPhoto(photoid int) (bool, error) {
	stmt := `
  	MATCH (p:Photo)<-[r:UPLOAD]-(u:User) WHERE ID(p) = {photoid}
    DETACH DELETE p
  	`
	params := map[string]interface{}{"photoid": photoid}
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

// CheckExistUserPhoto func
func CheckExistUserPhoto(photoid int) (bool, error) {
	return CheckExistNodeWithID(photoid)
}

// GetUserIDUploadedPhoto func
func GetUserIDUploadedPhoto(photoid int) (int, error) {
	stmt := `
    MATCH (u:User)-[r:UPLOAD]->(p:Photo) WHERE ID(p) = {photoid}
    RETURN ID(u) AS id
  	`
	params := map[string]interface{}{"photoid": photoid}
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

// GetUserPhoto func
func GetUserPhoto(photoid int, myuserid int) (models.Photo, error) {
	stmt := `
		MATCH(me:User) WHERE ID(me) = {myuserid}
		MATCH (p:Photo)<-[:UPLOAD]-(u:User)
		WHERE ID(p) = {photoid} AND ( p.privacy = 1 OR (p.privacy = 2 AND exists(me-[:FOLLOW]->(u))) OR u = me)
		RETURN
			ID(s) AS id, p.message AS message, p.photo AS photo, p.created_at AS created_at, p.updated_at AS updated_at,
			case p.privacy when null then 1 else p.privacy end AS privacy, case p.status when null then 1 else p.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			s.likes AS likes, p.comments AS comments, p.shares AS shares,
			exists((me)-[:LIKE]->(s)) AS is_liked
		`
	params := map[string]interface{}{"photoid": photoid, "myuserid": myuserid}
	res := []models.Photo{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Photo{}, err
	}
	if len(res) > 0 && res[0].PhotoID == photoid {
		return res[0], nil
	}
	return models.Photo{}, nil
}

// CreatePhotoLike func
func CreatePhotoLike(photoid int, userid int) (bool, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u) = {userid}
	MATCH(p:Photo) WHERE ID(p) = {photoid}
	MERGE(u)-[l:LIKE]->(p)
	ON CREATE SET l.created_at = TIMESTAMP()
	RETURN exists((u)-[l]->(s)) AS liked
	`
	params := map[string]interface{}{"photoid": photoid, "userid": userid}
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

// GetPhotoLikes func
func GetPhotoLikes(photoid int, myuserid int, orderby string, skip int, limit int) ([]models.SUserLike, error) {

	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {myuserid}
	MATCH (u:User)-[l:LIKE]->(p:Photo)
	WHERE ID(p) = {photoid}
	RETURN
		ID(u) AS id, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
		l.created_at as liked_at,
		exists((me)-[:FOLLOW]->(u)) AS is_followed
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, orderby)
	params := map[string]interface{}{
		"photoid":  photoid,
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

// DeletePhotoLike func
func DeletePhotoLike(photoid int, userid int) (bool, error) {
	stmt := `
	MATCH (u:User)-[l:LIKE]->(p:Photo)
	WHERE ID(p) = {photoid} AND ID(u) = {userid}
	DELETE l
	`
	params := map[string]interface{}{"photoid": photoid, "userid": userid}
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

// CheckExistPhotoLike func
func CheckExistPhotoLike(photoid int, userid int) (bool, error) {
	stmt := `
  	MATCH (u:User)-[l:LIKE]->(p:Photo)
		WHERE ID(u) = {userid} AND ID(p) = {photoid}
		RETURN ID(l) as id
  	`
	params := neoism.Props{"photoid": photoid, "userid": userid}
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

// IncreasePhotoLikes func
func IncreasePhotoLikes(photoid int) (bool, error) {
	stmt := `
	MATCH (p:Photo)
	WHERE ID(p)= {photoid}
	SET p.likes = p.likes+1
	RETURN ID(p) AS id
	`
	params := neoism.Props{"photoid": photoid}
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
	if len(res) > 0 && res[0].ID == photoid {
		return true, nil
	}
	return false, nil
}

// DecreasePhotoLikes func
func DecreasePhotoLikes(photoid int) (bool, error) {
	stmt := `
	MATCH (p:Photo)
	WHERE ID(p)= {photoid}
	SET p.likes = p.likes-1
	RETURN ID(p) AS id
	`
	params := neoism.Props{"photoid": photoid}
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
	if len(res) > 0 && res[0].ID == photoid {
		return true, nil
	}
	return false, nil
}

// IncreasePhotoComments func
func IncreasePhotoComments(photoid int) (bool, error) {
	stmt := `
	MATCH (p:Photo)
	WHERE ID(p)= {photoid}
	SET p.comments = p.comments+1
	RETURN ID(p) AS id
	`
	params := neoism.Props{"photoid": photoid}
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
	if len(res) > 0 && res[0].ID == photoid {
		return true, nil
	}
	return false, nil
}

// DecreasePhotoComments func
func DecreasePhotoComments(photoid int) (bool, error) {
	stmt := `
	MATCH (p:Photo)
	WHERE ID(p)= {photoid}
	SET p.comments = p.comments-1
	RETURN ID(p) AS id
	`
	params := neoism.Props{"photoid": photoid}
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
	if len(res) > 0 && res[0].ID == photoid {
		return true, nil
	}
	return false, nil
}

// CheckPhotoInteractivePermission func to get a photo
func CheckPhotoInteractivePermission(photoid int, userid int) (bool, error) {
	stmt := `
		MATCH (who:User) WHERE ID(who) = {userid}
		MATCH (u:User)-[r:UPLOAD]->(p:Photo)
		WHERE ID(p) = {photoid}
		RETURN exist((who)-[:FOLLOW]->(u)) AS followed, p.privacy AS privacy, u = who AS owner
		`
	params := map[string]interface{}{"userid": userid, "photoid": photoid}
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
		if res[0].Privacy == configs.Public || (res[0].Followed && res[0].Privacy == configs.ShareToFollowers) || res[0].Owner {
			return true, nil
		}
	}
	return false, nil

}
