# Connect

API for Connect DAPP.

## Getting Started

```bash
brew install golang dep
go get github.com/MartinResearchSociety/connect
cd $GOPATH/src/github.com/MartinResearchSociety/connect
dep ensure
go run app/main.go
```

Swagger UI available at `localhost:8030/swagger`

## Developer Notes

- No views, only an optional swagger UI
- The only 2 unprotected endpoints are `/v1/user/signup` & `/v1/user/login`

### Sign Up (Create User)

Required to fetch token later.

```bash
curl -X POST "http://localhost:8030/v1/user/signup" -H "accept: application/json" -H "content-type: application/json" -d "{ \"email\": \"string@example.com\", \"first_name\": \"string\", \"id\": \"507f1f77bcf86cd799439011\", \"image\": \"string\", \"last_name\": \"string\", \"password\": \"string12345\", \"social_media\": [ { \"key\": \"string\", \"name\": \"string\", \"url\": \"string.com\" } ], \"username\": \"string\"}"
```

- If using swagger UI, note that you have to change ID to a BSON ObjectID
- The ID here doesn't matter, it just has to exist for validation
- In the DB, the ObjectID is overwritten to another ObjectID
- Should this field be removed from the request body and have the ObjectID created on the server side?

### Log In (Fetch Token)

```bash
curl -X POST "http://localhost:8030/v1/user/login" -H "accept: application/json" -H "content-type: application/json" -d "{ \"password\": \"string12345\", \"email\": \"string@example.com\"}"
```

- If using swagger UI, note that you have to replace `"username": "string"` with `"email": "string@example.com"`

### Get All Users

This is an example of a route that requires a token.

```bash
curl -X GET "http://localhost:8030/v1/user/" -H "accept: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MzIwMTMzOTcsImp0aSI6ImQ0MWQ4Y2Q5OGYwMGIyMDRlOTgwMDk5OGVjZjg0MjdlIiwiaWF0IjoxNTMxOTk5Nzk3LCJpc3MiOiJsb2NhbGhvc3QiLCJuYmYiOjE1MzE5OTk3OTcsInN1YiI6IjViNTA3NDc5NDk1YjNiMzUyZjNhNDNmNCJ9.V9vamU5S0_iVuk3KtkXupBDUhGVI3eZHRUm9k8zQ8PQ"
```
