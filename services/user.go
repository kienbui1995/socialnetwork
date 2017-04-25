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

	node := neoism.Node{}
	if len(user.About) > 0 {
		node.SetProperty("about", user.About)
	}
	if len(user.Avatar) > 0 {
		node.SetProperty("avatar", user.Avatar)
	}
	if len(user.Birthday) > 0 {
		node.SetProperty("birthday", user.Birthday)
	}
	if len(user.Cover) > 0 {
		node.SetProperty("cover", user.Cover)
	}
	if len(user.Email) > 0 {
		node.SetProperty("email", user.Email)
	}
	if len(user.FacebookID) > 0 {
		node.SetProperty("facebook_id", user.FacebookID)
	}
	if len(user.FacebookToken) > 0 {
		node.SetProperty("facebook_token", user.FacebookToken)
	}
	if len(user.FirstName) > 0 {
		node.SetProperty("first_name", user.FirstName)
	}
	if len(user.FullName) > 0 {
		node.SetProperty("full_name", user.FullName)
	}
	if len(user.Gender) > 0 {
		node.SetProperty("gender", user.Gender)
	}
	if len(user.LastName) > 0 {
		node.SetProperty("last_name", user.LastName)
	}
	if len(user.MiddleName) > 0 {
		node.SetProperty("middle_name", user.MiddleName)
	}
	type resStruct struct {
		ID int `json:"ID"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err == nil {
		return res[0].ID, nil
	}
	return 0, err
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
	MATCH (u:User) WHERE ID(u) = {userId} RETURN u.avatar as Avatar, ID(u) as ID, u.username as Username, u.password as Password, u.email as email, u.status as Status LIMIT 25;
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

	if len(res) != 0 {
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
	if mailactive.UserID != 0 && len(mailactive.ActiveCode) != 0 && len(mailactive.Email) != 0 {
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
		 OPTIONAL MATCH(u:User) WHERE u.username CONTAINS {s}  OR u.full_name  CONTAINS {s}
		 RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
		 exists((:User{username:{userid}})-[:FOLLOW]->(u)) as is_followed
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
