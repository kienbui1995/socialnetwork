package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
)

//neo := connector.Neoism{IP:"155.94.144.150", Port: 7474, User: "neo4j", Password: "tlis2016", Type: "http"}

// URLDB CONST
const URLDB = "http://neo4j:madawg00@localhost:7474/db/data/"

var conn, _ = neoism.Connect(URLDB)

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
