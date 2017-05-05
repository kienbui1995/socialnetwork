package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

//CreateUserSubscriber func
func CreateUserSubscriber(fromid int, toid int) (int, error) {
	stmt := `
	MATCH(from:User) WHERE ID(from) = {fromid}
  MATCH (to:User) WHERE ID(to) = {toid}
  CREATE UNIQUE (from)-[f:FOLLOW]->(to)
 	return ID(f) as id
	`
	res := []struct {
		ID int `json:"id"`
	}{}
	params := neoism.Props{
		"fromid": fromid,
		"toid":   toid,
	}
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
	return -1, errors.New("Don't create follow relationship")

}

//DeleteUserSubscriber fun
func DeleteUserSubscriber(fromid int, toid int) (bool, error) {
	stmt := `
  	MATCH (from:User)-[f:FOLLOW]->(to:User) WHERE ID(from) = {fromid} AND ID(to) = {toid} delete f
  	`
	params := neoism.Props{"fromid": fromid, "toid": toid}
	res := -1
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

//CheckExistUserSubscriber fun
func CheckExistUserSubscriber(fromid int, toid int) (bool, error) {
	stmt := `
  	MATCH (from:User)-[f:FOLLOW]->(to:User) WHERE ID(from) = {fromid} AND ID(to) = {toid} return ID(f) as id
  	`
	params := neoism.Props{"fromid": fromid, "toid": toid}
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
	if len(res) > 0 {
		if res[0].ID > 0 {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

// GetSubscribers func
func GetSubscribers(userid int) ([]models.SUser, error) {
	if (userid) <= 0 {
		return nil, fmt.Errorf("Error in get Subscribers services: userid is null: %v", userid)
	}

	stmt := `
	MATCH (userid:User)-[f:FOLLOW]->(u:User)
	WHERE ID(userid)= {userid}
	RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
	TRUE	AS is_followed
  	`
	params := neoism.Props{"userid": userid}
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
	if len(res) > 0 {
		if res[0].UserID > 0 {
			return res, nil
		}
		return nil, nil
	}
	return nil, nil

}

// GetFollowers func
func GetFollowers(userid int) ([]models.SUser, error) {
	if (userid) <= 0 {
		return nil, errors.New("Error in get Subscribers services: userid is null")
	}

	stmt := `
	MATCH (u:User)-[f:FOLLOW]->(user:User)
	WHERE ID(user)= {userid}
	RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
	exists((user)-[:FOLLOW]->(u)) as is_followed

  	`
	params := neoism.Props{"userid": userid}
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
	if len(res) > 0 {
		if res[0].UserID > 0 {
			return res, nil
		}

	}
	return nil, nil

}

// GetFollowerIDs func
func GetFollowerIDs(userid int) ([]int, error) {
	if (userid) <= 0 {
		return nil, errors.New("Error in get Subscribers services: userid is null")
	}

	stmt := `
	MATCH (u:User)-[f:FOLLOW]->(user:User)
	WHERE ID(user)= {userid}
	RETURN ID(u) as id

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
		return nil, err
	}
	if len(res) > 0 {
		var ids []int
		for index := 0; index < len(res); index++ {
			ids = append(ids, res[index].ID)
		}
		return ids, nil

	}
	return nil, nil

}

// IncreaseFollowersAndFollowings func
func IncreaseFollowersAndFollowings(fromid int, toid int) (bool, error) {
	if fromid < 0 || toid < 0 {
		return false, errors.New("Error in IncreaseFollowersAndFollowings services: fromid or toid is null")
	}
	stmt := `MATCH (u1:User),(u2:User)
	WHERE ID(u1)= {fromid} and ID(u2) = {toid}
	SET u1.followings = u1.followings+1, u2.followers = u2.followers+1 RETURN ID(u1) as id1, ID(u2) as id2
	`
	params := neoism.Props{"fromid": fromid, "toid": toid}
	res := []struct {
		ID1 int `json:"id1"`
		ID2 int `json:"id2"`
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
	if len(res) > 0 && res[0].ID1 == fromid && res[0].ID2 == toid {
		return true, nil
	}
	return false, nil
}

// DecreaseFollowersAndFollowings func
func DecreaseFollowersAndFollowings(fromid int, toid int) (bool, error) {
	if fromid <= 0 || toid <= 0 {
		return false, errors.New("Error in DecreaseFollowersAndFollowings services: fromid or toid is null")
	}
	stmt := `MATCH (u1:User),(u2:User)
	WHERE ID(u1)= {fromid} and ID(u2) = {toid}
	SET u1.followings = u1.followings-1, u2.followers = u2.followers-1 RETURN ID(u1) as id1, ID(u2) as id2
	`
	params := neoism.Props{"fromid": fromid, "toid": toid}
	res := []struct {
		ID1 int `json:"id1"`
		ID2 int `json:"id2"`
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
	if len(res) > 0 && res[0].ID1 == fromid && res[0].ID2 == toid {
		return true, nil
	}
	return false, nil
}
