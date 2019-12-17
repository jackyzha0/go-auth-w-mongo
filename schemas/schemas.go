package schemas

type User struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
