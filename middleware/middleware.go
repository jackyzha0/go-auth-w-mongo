// Package middleware defines possible middleware
// that can be used by the router.
package middleware

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/jackyzha0/go-auth-w-mongo/db"
	"github.com/jackyzha0/go-auth-w-mongo/schemas"
)

// Auth is a simple authentication middleware with an option to
// check for if the user is an admin or not
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
			log.Warn("Bad Auth Attempt: Could not read cookie.")
			return
		}

		// if no err, get cookie value
		sessionToken := c.Value

		filter := bson.M{"sessionToken": sessionToken}
		var res schemas.User
		findErr := db.Users.Find(filter).One(&res)

		if findErr != nil {

			// no user with matching session_token
			if findErr == mgo.ErrNotFound {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				log.Warnf("Bad Auth Attempt: No user with token %s.", sessionToken)
				return
			}

			// other error
			w.WriteHeader(http.StatusInternalServerError)
			log.Warn("Bad Auth Attempt: Internal Error when finding user.")
			return
		}

		// parse time
		expireTime, timeParseErr := time.Parse(time.RFC3339, res.SessionExpires)

		// token time invalid
		if timeParseErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Warn("Bad Auth Attempt: Session expiry date wrong.")
			return
		}

		// token expired
		if time.Now().After(expireTime) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if adminCheck && !res.IsAdmin {
			w.WriteHeader(http.StatusUnauthorized)
			log.Warn("Bad Auth Attempt: Not admin. Attempt from user %v", res.Email)
			return
		}

		// token ok
		r.Header.Set("X-res-email", res.Email)
		req(w, r)
	}
}
