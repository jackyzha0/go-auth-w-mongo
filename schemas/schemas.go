package schemas

type User struct {
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
