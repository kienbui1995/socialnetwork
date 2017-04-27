package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/socialnetwork/configs"
)

//neo := connector.Neoism{IP:"155.94.144.150", Port: 7474, User: "neo4j", Password: "tlis2016", Type: "http"}

var conn, _ = neoism.Connect(configs.URLDB)

// CheckExistNode func to return quantity of nodes
func CheckExistNode(label string, where string) (bool, error) {
	stmt := `
	MATCH (u:{label}) WHERE {where} RETURN ID(u) as id;
	`
	params := neoism.Props{"label": label, "where": where}

	res := []interface{}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	fmt.Printf("\nid: %v, %T", res, res)
	if err != nil {
		return false, err
	}

	if len(res) != 0 {

		return true, nil
	}
	return false, nil
}

// CheckExistNodeWithID func to return quantity of nodes
func CheckExistNodeWithID(id int) (bool, error) {
	stmt := `
	MATCH (u) WHERE ID(u)={id} RETURN ID(u) as id;
	`
	params := neoism.Props{"id": id}

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

	if len(res) != 0 {
		if res[0].ID == id {
			return true, nil
		}
	}
	return false, nil
}
