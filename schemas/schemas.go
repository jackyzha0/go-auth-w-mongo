package schemas

type User struct {
	Password       string `json:"password"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	SessionToken   string `json:"session_token"`
	SessionExpires string `json:"session_expires"`
}

type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
