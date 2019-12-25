[![GoDoc](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo?status.svg)](https://godoc.org/github.com/jackyzha0/go-auth-w-mongo)
[![GoReportCard](https://goreportcard.com/badge/github.com/jackyzha0/go-auth-w-mongo)](https://goreportcard.com/report/github.com/jackyzha0/go-auth-w-mongo)
# Go Auth App
### Simple session based authentication with Mux and MongoDB

This repository was created as an exercise in Go development. The following routes have been implemented:

```go
/register "register a new user (Admin Only)"
/login "creates a token if user is registered"
/dashboard "simple dashboard to display user's name"
```

You can build the binary by running `go build` and run it via `./go-auth-w-mongo` Run the sanity check tests by doing `go test`
