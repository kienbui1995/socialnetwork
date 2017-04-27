package configs

//Const config system
const (
	neo4jURL = "Bolt://neo4j:tlis2016@localhost:7687"
	URLDB    = "http://neo4j:madawg00@localhost:7474/db/data/"
	APIPort  = 8080
)

//PrivacyType uint
type PrivacyType uint

// Const privacy
const (
	Public PrivacyType = iota + 1
	ShareToFollowers
	Private
)
