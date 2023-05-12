package utils

// struct used to list all users
type ResponseUsers struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	PhoneNum  string `json:"phonenum"`
	Block     bool   `json:"block"`
}
