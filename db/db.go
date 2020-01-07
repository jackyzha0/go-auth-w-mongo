package db

import (
	"github.com/globalsign/mgo"
    "crypto/tls"
    "net"
)

var Session *mgo.Session
var Users *mgo.Collection

const (
	DB_LOCAL = 0
	DB_LOCAL_W_AUTH = 1
	DB_ATLAS = 2
)

const DB_TYPE = DB_LOCAL

// Use atlas connection
func init() {

	switch DB_TYPE {
	case DB_LOCAL:

		// session is a connection to the given URI
		Session, _ = mgo.Dial("mongodb://localhost:27017")
		// Users is a new connection to Users Collection
		Users = Session.DB("exampleDB").C("Users")

	case DB_LOCAL_W_AUTH:

		Session, _ = mgo.Dial("mongodb://localhost:27017")

		// set credentials
		cred := mgo.Credential{
			Username: "username",
			Password: "password",
		}

		// attempt to login
		err := Session.Login(&cred)
		if err != nil {
			panic(err)
		}

		Users = Session.DB("exampleDB").C("Users")

	case DB_ATLAS:

		// set Atlas URI 
		mongoURI := "mongodb://<username>:<password>@<shard of ip (should end in .mongodb.net)>:27017"
		dialInfo, err := mgo.ParseURL(mongoURI)

		tlsConfig := &tls.Config{}
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		    conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		    return conn, err
		}
		
		// dial session with info
		Session, err = mgo.DialWithInfo(dialInfo)

		if err != nil {
			panic(err)
		}

		// Define Connections to Databases
		Users = Session.DB("alc_data").C("Users")
	}
}
