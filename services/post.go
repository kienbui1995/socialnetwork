package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/configs"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUserPost func
func CreateUserPost(userid int, message string, photo string, privacy int, status int) (int, error) {
	var p interface{}
	var stmt string
	if len(photo) == 0 {
		p = neoism.Props{
			"message":  message,
			"privacy":  privacy,
			"status":   status,
			"likes":    0,
			"comments": 0,
			"shares":   0,
		}
		stmt = `
	    MATCH(u:User) WHERE ID(u) = {fromid}
	  	CREATE (s:Status:Post { props } )<-[r:POST]-(u)
			SET s.created_at = TIMESTAMP()
			RETURN ID(s) as id
	  	`
	} else {
		p = neoism.Props{
			"message":  message,
			"photo":    photo,
			"privacy":  privacy,
			"status":   status,
			"likes":    0,
			"comments": 0,
			"shares":   0,
		}
		stmt = `
	    MATCH(u:User) WHERE ID(u) = {fromid}
	  	CREATE (s:Photo:Post { props } )<-[r:POST]-(u)
			SET s.created_at = TIMESTAMP()
			RETURN ID(s) as id
	  	`
	}
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

	return -1, nil
}

//GetUserPosts func
func GetUserPosts(userid int, myuserid int, orderby string, skip int, limit int, typepost int) ([]models.Post, error) {
	var stmt string
	if typepost == configs.PostPhoto {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Photo:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id, s.message AS message,
					case s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username,
					exists((me)-[:LIKE]->(s)) AS is_liked
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, orderby)
	} else if typepost == configs.PostStatus {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Status:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id, s.message AS message,
					case s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username,
					exists((me)-[:LIKE]->(s)) AS is_liked
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, orderby)
	} else if typepost == configs.Post {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id, s.message AS message,
					case s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					ID(u) AS userid, u.avatar AS avatar, u.full_name AS full_name, u.username AS username,
					exists((me)-[:LIKE]->(s)) AS is_liked
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, orderby)
	}
	params := map[string]interface{}{
		"userid":   userid,
		"myuserid": myuserid,
		"skip":     skip,
		"limit":    limit,
	}
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
	return nil, nil
}

// UpdateUserPost func
func UpdateUserPost(postid int, message string, photo string, privacy int, status int) (models.Post, error) {
	var stmt string
	var params map[string]interface{}
	if len(photo) > 0 {
		stmt = `
			MATCH (s:Post)<-[r:POST]-(u:User)
			WHERE ID(s) = {postid}
			SET s.message = {message}, s.photo = {photo}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}, s:Photo
			RETURN
				ID(s) AS id, s.message AS message, s.photo AS photo, s.created_at AS created_at, s.updated_at AS updated_at,
				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
				ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
				exists((u)-[:LIKE]->(s)) AS is_liked
			`
		params = map[string]interface{}{"postid": postid, "message": message, "photo": photo, "privacy": privacy, "status": status}
	} else {
		stmt = `
  	MATCH (s:Post)<-[r:POST]-(u:User)
    WHERE ID(s) = {postid}
		SET s.message = {message}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}
    RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.photo when null then "" else s.photo end AS photo,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			exists((u)-[:LIKE]->(s)) AS is_liked
  	`
		params = map[string]interface{}{"postid": postid, "message": message, "privacy": privacy, "status": status}
	}

	res := []models.Post{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Post{}, err
	}
	if len(res) > 0 && res[0].PostID >= 0 {
		return res[0], nil
	}

	return models.Post{}, errors.New("Dont' update user status")
}

// DeleteUserPost func
func DeleteUserPost(postid int) (bool, error) {
	stmt := `
  	MATCH (s:Post)<-[r:POST]-(u:User) WHERE ID(s) = {postid}
		DETACH DELETE s
  	`
	params := map[string]interface{}{"postid": postid}
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

// GetUserIDByPostID func
func GetUserIDByPostID(postid int) (int, error) {
	stmt := `
    MATCH (u:User)-[r:POST]->(s:Post)
		WHERE ID(s) = {postid}
		RETURN ID(u) AS id
  	`
	params := map[string]interface{}{"postid": postid}
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

// CheckExistUserPost func
func CheckExistUserPost(postid int) (bool, error) {
	stmt := `
	MATCH (u:Post) WHERE ID(u)={postid} RETURN ID(u) AS id;
	`
	params := neoism.Props{"postid": postid}

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

	if len(res) > 0 && res[0].ID == postid {
		return true, nil
	}
	return false, nil
}

// GetUserPost func
func GetUserPost(postid int, myuserid int) (models.Post, error) {
	stmt := `
		MATCH(me:User) WHERE ID(me) = {myuserid}
		MATCH (s:Post)<-[:POST]-(u:User)
		WHERE ID(s) = {postid}
		RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.photo when null then "" else s.photo end AS photo,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			s.likes AS likes, s.comments AS comments, s.shares AS shares,
			exists((me)-[:LIKE]->(s)) AS is_liked
		`
	params := map[string]interface{}{"postid": postid, "myuserid": myuserid}
	res := []models.Post{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Post{}, err
	}
	if len(res) > 0 && res[0].PostID == postid {
		return res[0], nil

	}

	return models.Post{}, errors.New("ERROR in GetUserPost service: Postid <0")
}

// CreatePostLike func
func CreatePostLike(postid int, userid int) (int, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u) = {userid}
	MATCH(s:Post) WHERE ID(s) = {postid}
	MERGE(u)-[l:LIKE]->(s)
	ON CREATE SET l.created_at = TIMESTAMP()
	RETURN exists((u)-[l]->(s)) AS liked, s.likes AS likes
	`
	params := map[string]interface{}{"postid": postid, "userid": userid}
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

// GetPostLikes func
func GetPostLikes(postid int, myuserid int, orderby string, skip int, limit int) ([]models.SUserLike, error) {

	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {myuserid}
	MATCH (u:User)-[l:LIKE]->(s:Post)
	WHERE ID(s) = {postid}
	RETURN
		ID(u) AS id, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
		l.created_at as liked_at,
		exists((me)-[:FOLLOW]->(u)) AS is_followed
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, orderby)
	params := map[string]interface{}{
		"postid":   postid,
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

// DeletePostLike func
func DeletePostLike(postid int, userid int) (int, error) {
	stmt := `
	MATCH (u:User)-[l:LIKE]->(s:Post) WHERE ID(s) = {postid} AND ID(u) = {userid}
	DELETE l
	RETURN s.likes AS likes
	`
	params := map[string]interface{}{"postid": postid, "userid": userid}
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

// CheckExistPostLike func
func CheckExistPostLike(postid int, userid int) (bool, error) {
	stmt := `
  	MATCH (u:User)-[l:LIKE]->(s:Post)
		WHERE ID(u) = {userid} AND ID(s) = {postid}
		RETURN ID(l) as id
  	`
	params := neoism.Props{"postid": postid, "userid": userid}
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

// IncreasePostLikes func
func IncreasePostLikes(postid int) (bool, error) {
	stmt := `
	MATCH (s:Post)
	WHERE ID(s)= {postid}
	SET s.likes = s.likes+1
	RETURN ID(s) AS id
	`
	params := neoism.Props{"postid": postid}
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
	if len(res) > 0 && res[0].ID == postid {
		return true, nil
	}
	return false, nil
}

// DecreasePostLikes func
func DecreasePostLikes(postid int) (bool, error) {
	stmt := `
	MATCH (s:Post)
	WHERE ID(s)= {postid}
	SET s.likes = s.likes-1
	RETURN ID(s) AS id
	`
	params := neoism.Props{"postid": postid}
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
	if len(res) > 0 && res[0].ID == postid {
		return true, nil
	}
	return false, nil
}

// CheckPostInteractivePermission func to check interactive permisson for user with a post
func CheckPostInteractivePermission(postid int, userid int) (bool, error) {
	stmt := `
		MATCH (who:User) WHERE ID(who) = {userid}
		MATCH (u:User)-[r:POST]->(s:Post)
		WHERE ID(s) = {postid}
		RETURN exists((who)-[:FOLLOW]->(u)) AS followed, s.privacy AS privacy
		`
	params := map[string]interface{}{"userid": userid, "postid": postid}
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
