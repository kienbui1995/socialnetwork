package services

import (
	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/models"
)

// CreateStatus func
func CreateStatus(post models.Post) (int, error) {
	p := neoism.Props{
		"message": post.Message,
		"photo":   post.Photo,
	}
	stmt := `
	CREATE (s:Status:Post { props } ) SET s.created_at = TIMESTAMP() RETURN ID(s) as id
	`
	params := map[string]interface{}{"props": p}
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

// GetFeedByUserID func

// // GetAllPost func
// func GetAllPost() ([]models.Post, error) {
// 	post := models.Post{}
// 	stmt := `
// 	MATCH (p:Post) RETURN ID(p), p.Content, p.Image, p.CreatedTime, p.UpdateTime, p.Status LIMIT 25;
// 	`
// 	res := []struct {
// 		//id       int     `json:ID(u)`
// 		PostID      int    `json:"ID(p)"`
// 		Content     string `json:"p.Content"`
// 		Image       string `json:"p.Image"`
// 		CreatedTime string `json:"p.CreatedTime"`
// 		UpdatedTime string `json:"p.UpdatedTime"`
// 		Status      int    `json:"u.Status"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement: stmt,
//
// 		Result: &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	k := int(len(res))
// 	fmt.Printf("%d\n", k)
// 	if k < 1 {
// 		return nil, errors.New("No Post")
// 	}
// 	var listpost []models.Post
// 	// fmt.Printf("%d\n", k)
// 	for i := range res {
//
// 		post.PostID = int(res[i].PostID)
// 		post.Content = res[i].Content
// 		post.Image = res[i].Image
// 		post.CreatedTime = res[i].CreatedTime
// 		post.UpdatedTime = res[i].UpdatedTime
// 		post.Status = int(res[i].Status)
//
// 		listpost = append(listpost, post)
// 		fmt.Printf("%d\n", i)
// 	}
// 	return listpost, nil
// }
//
// // GetPost func
// func GetPost(postid int) (models.Post, error) {
// 	var post = models.Post{}
// 	stmt := `
// 	MATCH (p:Post) WHERE ID(p) = {postid} RETURN ID(p), p.Content, p.Image , p.CreatedTime, p.UpdatedTime, p.Status LIMIT 25;
// 	`
// 	params := neoism.Props{"postid": postid}
// 	res := []struct {
// 		//id       int     `json:ID(u)`
// 		PostID      int    `json:"ID(p)"`
// 		Content     string `json:"p.Content"`
// 		Image       string `json:"p.Image"`
// 		CreatedTime string `json:"p.CreatedTime"`
// 		UpdatedTime string `json:"p.UpdatedTime"`
// 		Status      int    `json:"p.Status"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	// k := len(res)
// 	// // fmt.Printf("%d\n", k)
// 	// for index := 0; index < k; index++ {
// 	// 	fmt.Printf("id: %d,\nuname: %s,\npass: %s,\nemail: %s,\nstatus: %d\n",
// 	// 		int(res[index].UserID), res[index].Username, res[index].Password, res[index].Email, int(res[index].Status))
// 	// }
// 	if err != nil {
// 		return post, err
// 	}
// 	if len(res) == 1 {
// 		post = models.Post{PostID: res[0].PostID, Content: res[0].Content, Image: res[0].Image, CreatedTime: res[0].CreatedTime, UpdatedTime: res[0].UpdatedTime, Status: res[0].Status}
// 		return post, nil
// 	} else if len(res) > 1 {
// 		return post, errors.New("Many User")
// 	} else {
// 		return post, errors.New("No User")
// 	}
// }
//
// // UpdatePost func
// func UpdatePost(post models.Post) (bool, error) {
//
// 	stmt := `
// 	MATCH (p:Post) WHERE ID(p) = {postid} SET p.Content = {content}, p.Image = {image}, p.Status = {status}, p.UpdatedTime = {updatetime};
// 	`
// 	params := neoism.Props{"postid": post.PostID, "content": post.Content, "image": post.Image, "status": post.Status, "updatetime": post.UpdatedTime}
// 	res := false
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
//
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }
//
// // DeletePost func
// func DeletePost(postid int) (int, error) {
// 	stmt := `
// 	MATCH (p:Post) WHERE ID(p) = {postid} delete p RETURN count(p) ;
// 	`
// 	params := neoism.Props{"postid": postid}
// 	res := -1
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
//
// 	err := conn.Cypher(&cq)
// 	return res, err
// }
//
// //CheckExistPost func to check exist User
// func CheckExistPost(postid int) (bool, error) {
// 	where := fmt.Sprintf("ID(p) = %d", postid)
// 	existNode, err := CheckExistNode("Post", where)
// 	if err != nil {
// 		return false, err
// 	}
// 	if existNode == true {
// 		return true, nil
// 	}
// 	return false, nil
// }
//
// //GetPostByUserID func to get post info by userid of who write it
// func GetPostByUserID(userid int) ([]models.Post, error) {
// 	var post models.Post
// 	stmt := `
// 	MATCH (u:User{ID(u) = {userid}})-[w:Write]->(p:Post) RETURN ID(p), p.Content, p.Image , p.CreatedTime, p.UpdatedTime, p.Status LIMIT 25;
// 	`
// 	params := neoism.Props{"userid": userid}
// 	res := []struct {
// 		//id       int     `json:ID(u)`
// 		PostID      int    `json:"ID(p)"`
// 		Content     string `json:"p.Content"`
// 		Image       string `json:"p.Image"`
// 		CreatedTime string `json:"p.CreatedTime"`
// 		UpdatedTime string `json:"p.UpdatedTime"`
// 		Status      int    `json:"p.Status"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	// k := len(res)
// 	// // fmt.Printf("%d\n", k)
// 	// for index := 0; index < k; index++ {
// 	// 	fmt.Printf("id: %d,\nuname: %s,\npass: %s,\nemail: %s,\nstatus: %d\n",
// 	// 		int(res[index].UserID), res[index].Username, res[index].Password, res[index].Email, int(res[index].Status))
// 	// }
// 	if err != nil {
// 		return nil, err
// 	}
// 	var listpost []models.Post
// 	for i := range res {
// 		post = models.Post{PostID: res[i].PostID, Content: res[i].Content, Image: res[i].Image, CreatedTime: res[i].CreatedTime, UpdatedTime: res[i].UpdatedTime, Status: res[i].Status}
// 		listpost = append(listpost, post)
// 	}
// 	return listpost, nil
// }
