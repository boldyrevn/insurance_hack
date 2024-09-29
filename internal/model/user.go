package model

const UserKey = "User"

type User struct {
	Login     string `json:"login,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Age       int    `json:"age,omitempty"`
}
