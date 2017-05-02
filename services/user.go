package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateUser func
func CreateUser(user models.User) (int, error) {
	stmt := `
	Create (u:User{
		username: {username},
		password: {password},
		email: {email},
		first_name: {firstname},
		middle_name: {middlename},
		last_name: {lastname},
		full_name: {fullname},
		about: {about},
		gender: {gender},
		birthday: {birthday},
		avatar: {avatar},
		cover: {cover},
		status: {status},
		is_vertified: {isvertified},
		updated_at: {updatedat},
		created_at: {createdat},
		facebook_id: {facebookid},
		facebook_token: {facebooktoken},
		posts: {posts},
		followers: {followers},
		followings: {followings}
		})
	return ID(u) as ID
	`
	params := neoism.Props{
		"username":      user.Username,
		"password":      user.Password,
		"email":         user.Email,
		"firstname":     user.FirstName,
		"middlename":    user.MiddleName,
		"lastname":      user.LastName,
		"fullname":      user.FullName,
		"about":         user.About,
		"gender":        user.Gender,
		"birthday":      user.Birthday,
		"avatar":        user.Avatar,
		"cover":         user.Cover,
		"status":        user.Status,
		"isvertified":   user.IsVertified,
		"updatedat":     user.UpdatedAt,
		"createdat":     user.CreatedAt,
		"facebookid":    user.FacebookID,
		"facebooktoken": user.FacebookToken,
		"posts":         user.Posts,
		"followers":     user.Followers,
		"followings":    user.Followings,
	}
	type resStruct struct {
		ID int `json:"ID"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err == nil {
		return res[0].ID, nil
	}
	return 0, err
}

// CreateUserTest func
func CreateUserTest(user models.User) (int, error) {
	p := neoism.Props{
		"username":     user.Account.Username,
		"password":     user.Account.Password,
		"email":        user.Account.Email,
		"is_vertified": user.Account.IsVertified,
	}

	if len(user.About) > 0 {
		p["about"] = user.About
	}
	if len(user.Avatar) > 0 {
		p["avatar"] = user.Avatar
	}
	if len(user.Birthday) > 0 {
		p["birthday"] = user.Birthday
	}
	if len(user.Cover) > 0 {
		p["cover"] = user.Cover
	}
	if len(user.FacebookID) > 0 {
		p["facebook_id"] = user.FacebookID
	}
	if len(user.FacebookToken) > 0 {
		p["facebook_token"] = user.FacebookToken
	}
	if len(user.FirstName) > 0 {
		p["first_name"] = user.FirstName
	}
	if len(user.FullName) > 0 {
		p["full_name"] = user.FullName
	}
	if len(user.Gender) > 0 {
		p["gender"] = user.Gender
	}
	if len(user.LastName) > 0 {
		p["last_name"] = user.LastName
	}
	if len(user.MiddleName) > 0 {
		p["middle_name"] = user.MiddleName
	}
	p["posts"] = 0
	p["followers"] = 0
	p["followings"] = 0
	p["status"] = 0
	stmt := `
	CREATE (u:User { props } ) SET u.created_at = TIMESTAMP() RETURN ID(u) as id
	`
	res := []struct {
		ID int `json:"id"`
	}{}
	params := map[string]interface{}{"props": p}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	if err := conn.Cypher(&cq); err != nil {
		return -1, err
	}
	if len(res) > 0 {
		if res[0].ID >= 0 {
			return res[0].ID, nil
		}
	}
	return -1, nil
}

// GetAllUser func
func GetAllUser() ([]models.User, error) {
	return GetAllUserWithSkipLimit(0, 25)
}

// GetAllUserWithSkipLimit func
func GetAllUserWithSkipLimit(skip int, limit int) ([]models.User, error) {

	stmt := fmt.Sprintf("MATCH (u:User) RETURN ID(u) as ID, u.username as Username, u.password as Password, u.email as Email, u.status as Status, u.avatar as Avatar SKIP %v LIMIT %v",
		skip, limit)

	res := []models.User{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	k := int(len(res))

	if k < 1 {
		return nil, errors.New("No User")
	}

	return res, nil
}

// GetUser func
func GetUser(userid int) (models.User, error) {
	var user = models.User{}
	stmt := `
	MATCH (u:User)
	WHERE ID(u) = {userId}
	RETURN
	u.avatar as avatar, u.about as about, u.birthday as birthday, u.gender as gender, u.cover as cover,
	ID(u) as ID,
	u.username as username,
	u.full_name as full_name, u.first_name as first_name, u.last_name as last_name,
	u.email as email, u.status as Status,
	u.followers as followers, u.followings as followings, u.posts as posts,
	u.created_at as created_at, u.updated_at as updated_at
	LIMIT 25;
	`
	params := neoism.Props{"userId": userid}

	res := []models.User{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)

	if err != nil {
		return user, err
	}
	if len(res) == 1 {

		return res[0], nil
	} else if len(res) > 1 {
		return user, errors.New("Many User")
	} else {
		return user, errors.New("No User")
	}
}

// UpdateUser func
func UpdateUser(user models.User) (bool, error) {

	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} SET u.Email = {email}, u.Status = {status},  u.Avatar = {avatar}, u.Password = {password};
	`
	params := neoism.Props{"userid": user.UserID, "email": user.Email, "status": user.Status, "avatar": user.Avatar, "password": user.Password}
	res := false
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

// DeleteUser func
func DeleteUser(userid int) (int, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid} delete u RETURN count(u) ;
	`
	params := neoism.Props{"userid": userid}
	res := -1
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	return res, err
}

//CheckExistUserByID func to check exist User
func CheckExistUserByID(userid int) (bool, error) {
	where := fmt.Sprintf("ID(u) = %d", userid)
	existNode, err := CheckExistNode("User", where)
	if err != nil {
		return false, err
	}
	if existNode == true {
		return true, nil
	}
	return false, nil
}

//CheckExistUserWithUsernameAndEmail func to check exist User
// return 1: Exist
// return 2: Exist Username
// return 3: Exist Email
// return 0: Don't exist
func CheckExistUserWithUsernameAndEmail(username string, email string) (int, error) {
	var where string
	if len(username) != 0 && len(email) != 0 {
		where = fmt.Sprintf("u.username = %s, u.email = %s", username, email)
		existNode, err := CheckExistNode("User", where)
		if err != nil {
			return 0, err
		}
		if existNode == true {
			return 1, nil
		}
	} else if len(username) != 0 {
		where = fmt.Sprintf("u.username = %s", username)
		existNode, err := CheckExistNode("User", where)
		if err != nil {
			return 0, err
		}
		if existNode == true {
			return 2, nil
		}
	} else if len(email) != 0 {
		where = fmt.Sprintf("u.email = %s", email)
		existNode, err := CheckExistNode("User", where)
		if err != nil {
			return 0, err
		}
		if existNode == true {
			return 3, nil
		}

	} else {
		return 0, errors.New("Missing username and email")
	}

	return 0, nil
}

//CheckExistUsername func
func CheckExistUsername(username string) (bool, error) {
	stmt := fmt.Sprintf("MATCH (u:User) WHERE u.username = \"%s\" RETURN ID(u) as id", username)
	type resStruct struct {
		ID int `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) > 0 && res[0].ID > 0 {
		return true, nil
	}
	return false, nil
}

//CheckExistEmail func
func CheckExistEmail(email string) (bool, error) {
	stmt := fmt.Sprintf("MATCH (u:User) WHERE u.email = \"%s\" RETURN ID(u) as id", email)
	type resStruct struct {
		ID int `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) != 0 {
		return true, nil
	}
	return false, nil
}

//CheckExistFacebookID func
func CheckExistFacebookID(facebookid string) (int, error) {
	stmt := fmt.Sprintf("MATCH (u:User) WHERE u.facebook_id = \"%s\" RETURN ID(u) as id", facebookid)
	type resStruct struct {
		ID int `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return 0, err
	}

	if len(res) != 0 {
		return res[0].ID, nil
	}
	return 0, nil
}

//CreateEmailActive func
func CreateEmailActive(mailactive models.EmailActive) error {
	if mailactive.UserID >= 0 && len(mailactive.ActiveCode) != 0 && len(mailactive.Email) != 0 {
		stmt := fmt.Sprintf("MATCH (u:User) WHERE ID(u) = \"%d\" AND u.email = \"%s\" SET u.active_code = \"%s\"", mailactive.UserID, mailactive.Email, mailactive.ActiveCode)
		cq := neoism.CypherQuery{
			Statement: stmt,
		}
		if err := conn.Cypher(&cq); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Error in creating email active")
}

//CreateRecoverPassword func
func CreateRecoverPassword(email string, recoverycode string) error {
	if len(email) != 0 && len(recoverycode) != 0 {
		stmt := fmt.Sprintf("MATCH (u:User) WHERE u.email = \"%s\" SET u.recovery_code = \"%s\", u.recovery_expired_at = TIMESTAMP()+1800000", email, recoverycode)
		cq := neoism.CypherQuery{
			Statement: stmt,
		}
		if err := conn.Cypher(&cq); err != nil {
			return err
		}
		return nil
	}
	return errors.New("Error in creating recovery code")
}

//VerifyRecoveryCode func
func VerifyRecoveryCode(email string, recoverycode string) (int, error) {
	if len(email) != 0 && len(recoverycode) != 0 {
		stmt := fmt.Sprintf("MATCH (u:User) WHERE u.email = \"%s\" and u.recovery_code = \"%s\"  and u.recovery_expired_at > TIMESTAMP() RETURN ID(u) as id", email, recoverycode) //needfix
		res := []struct {
			ID int `json:"id"`
		}{}
		cq := neoism.CypherQuery{
			Statement: stmt,
			Result:    &res,
		}
		if err := conn.Cypher(&cq); err != nil {
			return -1, err
		}
		if len(res) > 0 {
			return res[0].ID, nil
		}
		return -1, nil

	}
	return -1, errors.New("Error in verify recovery code")
}

//AddUserRecoveryKey func to add a property of user
func AddUserRecoveryKey(userid int, value interface{}) error {
	stmt := `
	MATCH(u:User) WHERE ID(u)= {userid} SET u.recovery_key = {value}
	`
	params := neoism.Props{
		"userid": userid,
		"value":  value,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}
	err := conn.Cypher(&cq)
	return err
}

//RenewPassword func
func RenewPassword(userid int, recoverykey string, newpassword string) (bool, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u)= {userid} AND u.recovery_key = {recoverykey} SET u.password = {newpassword}
	RETURN u.password as password
	`
	res := []struct {
		Password string `json:"password"`
	}{}
	params := neoism.Props{
		"userid":      userid,
		"recoverykey": recoverykey,
		"newpassword": newpassword,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	if len(res) > 0 && res[0].Password == newpassword {
		return true, nil
	}
	return false, nil
}

//DeleteRecoveryProperty func
func DeleteRecoveryProperty(userid int) (bool, error) {
	stmt := fmt.Sprintf("MATCH (u:User) WHERE ID(u) = %d REMOVE u.recovery_key, u.recovery_code, u.recovery_expired_at RETURN ID(u) as id, exists(u.recovery_key) as k, exists(u.recovery_code) as c,exists(u.recovery_expired_at) as e ", userid)
	type resStruct struct {
		ID int  `json:"id"`
		K  bool `json:"k"`
		C  bool `json:"c"`
		E  bool `json:"e"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) != 0 {
		if res[0].ID == userid && res[0].K == true && res[0].C == true && res[0].E == true {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("No exist user")
}

//FindUserByUsernameAndFullName func
func FindUserByUsernameAndFullName(userid int, s string) ([]models.SUser, error) {
	stmt := `
		 MATCH(a:User) where ID(a)={userid}
		 OPTIONAL MATCH(u:User) WHERE u.username CONTAINS {s}  OR u.full_name  CONTAINS {s}
		 RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
		 exists((a)-[:FOLLOW]->(u)) as is_followed
	`
	res := []models.SUser{}
	params := neoism.Props{
		"userid": userid,
		"s":      s,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		if len(res) == 1 && res[0].UserID == 0 {
			return nil, nil
		}
		return res, nil
	}
	return nil, nil
}

//GetNewsFeed func
func GetNewsFeed(userid int, orderby string, skip int, limit int) ([]models.UserStatus, error) {
	stmt := fmt.Sprintf(`
	MATCH(u:User) WHERE ID(u)= {userid}
	MATCH(u)-[:FOLLOW]->(u1:User)-[:POST]->(s:Status)
	RETURN
		ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
		case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
		ID(u1) AS userid, u1.full_name AS full_name, u1.avatar AS avatar, u1.username AS username,
		s.likes AS likes, s.comments AS comments, s.shares AS shares,
		exists((u)-[:LIKE]->(s)) AS is_liked
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, orderby)
	res := []models.UserStatus{}
	params := map[string]interface{}{
		"userid": userid,
		"skip":   skip,
		"limit":  limit,
	}
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

//IncreasePosts func
func IncreasePosts(userid int) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u)= {userid}
	SET u.posts = u.posts+1
	RETURN ID(u) AS id
	`
	params := neoism.Props{"userid": userid}
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
	if len(res) > 0 && res[0].ID == userid {
		return true, nil
	}
	return false, nil
}

//DecreasePosts func
func DecreasePosts(userid int) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u)= {userid}
	SET u.posts = u.posts-1
	RETURN ID(u) AS id
	`
	params := neoism.Props{"userid": userid}
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
	if len(res) > 0 && res[0].ID == userid {
		return true, nil
	}
	return false, nil
}
