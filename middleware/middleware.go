package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jackyzha0/go-auth-w-mongo/routes"
	"github.com/jackyzha0/go-auth-w-mongo/schemas"
)

func Auth(req http.HandlerFunc, adminCheck bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, cookieFetchErr := r.Cookie("session_token")

		// not auth'ed, redirect to login
		if cookieFetchErr != nil {
			if cookieFetchErr == http.ErrNoCookie {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Bad Request, could not read cookie.\n")
			return
		}

		// if no err, get cookie value
		sessionToken := c.Value

		filter := bson.M{"sessionToken": sessionToken}
		var res schemas.User
		findErr := routes.Users.Find(filter).One(&res)

		if findErr != nil {

			// no user with matching session_token
			if findErr == mgo.ErrNotFound {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				fmt.Fprintf(w, "Bad Request, no user with that email exists.\n")
				return
			}

			// other error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// parse time
		expireTime, timeParseErr := time.Parse(time.RFC3339, res.SessionExpires)

		// token time invalid
		if timeParseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// token expired
		if time.Now().After(expireTime) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if adminCheck && !res.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Bad Request, not admin.\n")
			return
		}

		// token ok
		r.Header.Set("X-res-email", res.Email)
		req(w, r)
	}
}
