[![GoDoc](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo?status.svg)](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo)
[![GoReportCard](https://goreportcard.com/badge/github.com/jackyzha0/go-auth-w-mongo)](https://goreportcard.com/report/github.com/jackyzha0/go-auth-w-mongo)
# Go Auth App
### Simple session based authentication with Mux and MongoDB

This repository was created as an exercise in Go development. This serves as a barebones boilerplate for a webapp that requires authentication via a database.

### Setup
You can build the binary by running `go build` and run it via `./go-auth-w-mongo`

Please make sure you have a local MongoDB instance running before attempting to run anything here! A simple auth flow may look something like this

1. User attempts to access dashboard
2. User is redirected to log in
3. User logs in and is redirected to dashboard
4. (Admin only) Admin can add a new user

The following routes have been implemented:

```go
/register "register a new user"
/login "authenticates a user"
/dashboard "simple dashboard to display user's name"
```

### `/register`
The `/register` endpoint first checks to see if the user has a valid session (this is done through the authentication middleware in `middleware/middleware.go`) and is an admin.

Then, it attempts to parse the form data into the User schema. After doing so, a hash is generated from the password to avoid storing it in plaintext, and the document is inserted into the database.

```bash
curl --location --request POST 'localhost:8080/register' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'email=test@email.com' \
--data-urlencode 'name=jacky' \
--data-urlencode 'password=pass123'
```

### `/login`
The `/login` endpoint attempts to decode the form data into the Credentials schema. Then, it attempts to find a user with a matching email in the database. If found, it compares the password hashes to see if the password is correect. If so, it will create a new session token, write that to both the client (as a cookie) and the database (as a document), and redirect them to `/dashboard`

```bash
curl --location --request POST 'localhost:8080/login' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'email=test@email.com' \
--data-urlencode 'password=pass123'
```

### `/dashboard`
The `/dashboard` endpoint takes advantage of the fact that our middleware injects a header called `X-res-email` with the user's email if the user has a valid session. This is a basic endpoint that returns a custom greeting based on the user's name.

```bash
curl --location --request GET 'localhost:8080/dashboard'
```

## Adding to the project
As this is just a boilerplate, this project is meant to be easily extensible. Feel free to add more endpoints in `routes/routes.go`, render and serve templates, and add more middleware! The project is your oyster :)

## FAQ

#### How does session based authentication work?

This photo (courtesy of PracticalDev) does a pretty good job at explaining it!

![Session Based Auth](https://res.cloudinary.com/practicaldev/image/fetch/s--jzM6Wq6e--/c_limit%2Cf_auto%2Cfl_progressive%2Cq_auto%2Cw_880/https://cdn-images-1.medium.com/max/800/0%2AP5OxJMihg0S0jyqk.png)

#### Adding an admin user
Remove the admin checker middleware from `server.go` register endpoint by changing line 22 from
```go
r.HandleFunc("/register", middleware.Auth(routes.Register, true))
```
to
```go
r.HandleFunc("/register", routes.Register)
```
Then, make a POST request to the `/register` endpoint with a `x-www-form-urlencoded` (form data) with the required fields. Make sure to set the admin field to true!


It would look something like this
```bash
curl --location --request POST 'localhost:8080/register' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'email=test@email.com' \
--data-urlencode 'name=jacky' \
--data-urlencode 'password=pass123' \
--data-urlencode 'admin=true'
```

Don't forget to change that line back later.

#### Changing URL of database

This is set in `routes/routes.go` line 19,

```go
var session, _ = mgo.Dial("mongodb://localhost:27017")
```

#### Changing name of database used and name of collection

This is set in `routes/routes.go` line 22,

```go
// "exampleDB" is the name of the database
// "Users" is the name of the collection
var Users = session.DB("exampleDB").C("Users")
```