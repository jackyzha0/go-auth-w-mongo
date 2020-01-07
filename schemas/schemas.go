// Schema Definitions for MongoDB Documents
package schemas

// Describes app user
type User struct {
	Password       string `json:"password" bson:"password"`
	Name           string `json:"name" bson:"name"`
	Email          string `json:"email" bson:"email"`
	SessionExpires string `json:"sessionExpires" bson:"sessionExpires"`
	SessionToken   string `json:"sessionToken" bson:"sessionToken"`
	IsAdmin        bool   `json:"isAdmin" bson:"isAdmin"`
}

// Describes login credentials
type Credentials struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}
