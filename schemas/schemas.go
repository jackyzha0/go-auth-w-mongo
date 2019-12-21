// Schema Definitions for MongoDB Documents
package schemas

// Describes app user
type User struct {
	Password       string `json:"password"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	SessionExpires string `json:"session_expires"`
	SessionToken   string `json:"session_token"`
}

// Describes login credentials
type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}