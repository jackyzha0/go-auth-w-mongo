[![GoDoc](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo?status.svg)](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo)
[![GoReportCard](https://goreportcard.com/badge/github.com/jackyzha0/go-auth-w-mongo)](https://goreportcard.com/badge/github.com/jackyzha0/go-auth-w-mongo)
# Go Auth App
### Simple session based authentication with Mux and MongoDB

This repository was created as an exercise in Go development. The following routes have been implemented:

```go
/register "register a new user"
/login "creates a token if user is registered"
/dashboard "simple dashboard to display user's name"
/dbhealthcheck "checks if connection to MongoDB is healthy"
```
