package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreatePost func
func CreatePost(post models.Post) (int, error) {
	p := neoism.Props{
		"postId":      int(post.PostID),
		"Message":     post.Message,
		"Type":        post.Type,
		"IsHidden":    post.IsHidden,
		"Status":      int(post.Status),
		"CreatedTime": post.CreatedTime,
		"UpdatedTime": post.UpdatedTime,
	}

	node, errNode := conn.CreateNode(p)
	if errNode != nil {
		return -1, errNode
	}
	errLabel := node.SetLabels([]string{"Post"})
	if errLabel != nil {
		node.Delete()
		return -1, errLabel
	}
	var propNode neoism.Props
	var err error
	propNode, err = node.Properties()
	if err != nil {
		return -1, err
	}

	postid := propNode["postId"].(float64)

	//create relate with user
	id := int(postid)
	return id, nil
}

// GetPost func
func GetPost(postid int) (models.Post, error) {
	var post = models.Post{}
	stmt := `
	MATCH (p:Post) WHERE p.postId = {postid} RETURN ID(p), p.Message, p.Type , p.CreatedTime, p.Status LIMIT 25;
	`
	params := neoism.Props{"postid": postid}
	res := []struct {
		//id       int     `json:ID(u)`
		PostID      int    `json:"ID(p)"`
		Message     string `json:"p.Message"`
		Type        string `json:"p.Type"`
		CreatedTime string `json:"p.CreatedTime"`
		Status      int    `json:"p.Status"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	// k := len(res)
	// // fmt.Printf("%d\n", k)
	// for index := 0; index < k; index++ {
	// 	fmt.Printf("id: %d,\nuname: %s,\npass: %s,\nemail: %s,\nstatus: %d\n",
	// 		int(res[index].UserID), res[index].Username, res[index].Password, res[index].Email, int(res[index].Status))
	// }
	if err != nil {
		return post, err
	}
	if len(res) == 1 {
		post = models.Post{PostID: res[0].PostID, Message: res[0].Message, Type: res[0].Type, CreatedTime: res[0].CreatedTime, Status: res[0].Status}
		return post, nil
	} else if len(res) > 1 {
		return post, errors.New("Many User!")
	} else {
		return post, errors.New("No User!")
	}
}

// UpdatePost func
func UpdatePost(user models.User) (bool, error) {

	stmt := `
	MATCH (u:User) WHERE u.userId = {userid} SET u.Email = {email}, u.Status = {status};
	`
	params := neoism.Props{"userid": user.UserID, "email": user.Email, "status": user.Status}
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

// DeletePost func
func DeletePost(userid int) (int, error) {
	stmt := `
	MATCH (u:User) WHERE u.userId = {userid} delete u RETURN count(u) ;
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

//CheckExistPost func to check exist User
func CheckExistPost(userid int) (bool, error) {
	where := fmt.Sprintf("u.userId = %d", userid)
	number, err := CheckExistNode("User", where)
	if err != nil {
		return false, err
	}
	if number > 0 {
		return true, nil
	}
	return false, nil
}
