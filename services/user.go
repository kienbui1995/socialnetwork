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

// GetAllUser func
func GetAllUser() ([]models.User, error) {

	stmt := `
	MATCH (u:User) RETURN ID(u) as UserID, u.Username as Username, u.Password as Password, u.Email as Email, u.Status as Status, u.Avatar as Avatar LIMIT 25;
	`
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

// GetAllUserWithSkipLimit func
func GetAllUserWithSkipLimit(skip int, limit int) ([]models.User, error) {

	stmt := fmt.Sprintf("MATCH (u:User) RETURN ID(u) as UserID, u.Username as Username, u.Password as Password, u.Email as Email, u.Status as Status, u.Avatar as Avatar SKIP %v LIMIT %v",
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
	MATCH (u:User) WHERE ID(u) = {userId} RETURN u.Avatar as Avatar, ID(u) as UserID, u.Username as Username, u.Password as Password, u.Email as email, u.Status as Status LIMIT 25;
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

//CheckExistUser func to check exist User
func CheckExistUser(userid int) (bool, error) {
	where := fmt.Sprintf("ID(u) = %d", userid)
	number, err := CheckExistNode("User", where)
	if err != nil {
		return false, err
	}
	if number > 0 {
		return true, nil
	}
	return false, nil
}
