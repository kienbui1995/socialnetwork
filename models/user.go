package models

type User struct {
	id       int
	username string
	password string
	email    string
	status   int
}

func testUser() string {
	return "kienbui"
}
