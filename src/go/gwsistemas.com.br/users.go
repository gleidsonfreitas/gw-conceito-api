package main

// GOOS=linux GOARCH=amd64 go build -o main users.go && zip users.zip main
// DB_DATABASE=development DB_HOST=gw.cp02vhkg15jt.us-east-1.rds.amazonaws.com DB_PASSWORD=gw.s1st3m4s DB_PORT=5432 DB_USER=postgres

import (
	"database/sql"
	"fmt"
	"strconv"

	"gwsistemas.com.br/helpers"
	"gwsistemas.com.br/models"
	"gwsistemas.com.br/repository"

	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
	fmt.Println(" + Users")
}

func main() {
	// HandleLambdaEvent(models.MyEvent{Method: "get-users", Params: []string{}, Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjEsImV4cCI6MTU3NjU0NTMwNSwiaXNzIjoibmF6YXIuaW8ifQ.avXj3QqQgaMh9an6JIpu0cAmLRkdDDyV0PBhMJy7o_M"})
	lambda.Start(HandleLambdaEvent)
}

func HandleLambdaEvent(event models.MyEvent) (models.MyResponse, error) {
	myResponse := models.MyResponse{}

	// ValidateToken
	token := event.Token
	claims, err := models.ValidateToken(token)
	if err != nil || claims.UserId == 0 {
		return models.UnauthorizedResponse(), nil
	}

	if event.Method == "create-user" {
		return CreateUser(event, claims), nil
	}

	if event.Method == "get-user" {
		return GetUser(event, claims), nil
	}

	if event.Method == "get-users" {
		return GetUsers(event, claims), nil
	}

	if event.Method == "authenticate-user" {
		return AuthenticateUser(event, claims), nil
	}

	return myResponse, nil
}

func Authorize(db *sql.DB, userId int) bool {
	// Current user
	user, err := repository.FindUser(db, userId)
	if err != nil {
		return false
	}

	if user.Enabled {
		return true
	}

	return false
}

/* --- User API --- */

// CreateUser creates a new user and stores it in the database
// if the email provided have not already be taken.
// It creates an api key and associates to an usage plan in api gateway.
// It also generates a JWT and stores the User.Id and User.Email inside.
//
// [POST /users/create]
func CreateUser(event models.MyEvent, claims models.MyClaims) models.MyResponse {
	// Params
	name := event.Params[0]
	email := event.Params[1]
	password := event.Params[2]

	// Connection
	db, err := helpers.Connection()
	if err != nil {
		return models.SomethingWentWrongResponse()
	}
	defer db.Close()

	// Authorized
	authorized := Authorize(db, claims.UserId)
	if !authorized {
		return models.UnauthorizedResponse()
	}

	// Email already used
	user, err := repository.FindUserByEmail(db, email)
	if err != nil {
		return models.SomethingWentWrongResponse()
	}
	if user != nil {
		return models.UserEmailAlreadyUsedResponse(user)
	}

	// Store the user in the database
	user, err = repository.CreateUser(db, name, email, password)
	if err != nil {
		return models.SomethingWentWrongResponse()
	}

	// Api key

	// JWT

	// User created
	return models.UserCreatedResponse(user)
}

// AuthenticateUser authenticates a user if the provided credentials matches
// the ones stored in the database.
// If the credentials are valid, it will refresh the JWT with a new expiration
// and store it in the database.
//
// [POST /users/authenticate]
func AuthenticateUser(event models.MyEvent, claims models.MyClaims) models.MyResponse {
	// Params
	email := event.Params[0]
	password := event.Params[1]

	// Connection
	db, err := helpers.Connection()
	if err != nil {
		return models.SomethingWentWrongResponse()
	}
	defer db.Close()

	// User not authenticated
	user, err := repository.FindUserByEmail(db, email)
	if err != nil {
		fmt.Println(err)
		return models.SomethingWentWrongResponse()
	}

	// User not authenticated
	if user.ValidatePassword(password) == false {
		fmt.Println(err)
		return models.SomethingWentWrongResponse()
	}

	// JWT
	user.Token = models.GenerateToken(user.Id)

	// User authenticated
	return models.UserAuthenticatedResponse(user)
}

// GetUser retrieves an user that matches the userId provided.
//
// [GET /users/{userId}]
func GetUser(event models.MyEvent, claims models.MyClaims) models.MyResponse {
	// Params
	userId, err := strconv.ParseInt(event.Params[0], 10, 64)
	if err != nil {
		return models.BadRequestResponse()
	}

	// Connection
	db, err := helpers.Connection()
	if err != nil {
		return models.SomethingWentWrongResponse()
	}
	defer db.Close()

	// Authorized
	authorized := Authorize(db, claims.UserId)
	if !authorized {
		return models.UnauthorizedResponse()
	}

	// User not found
	user, err := repository.FindUser(db, int(userId))
	if err != nil {
		return models.SomethingWentWrongResponse()
	}

	return models.UserFoundResponse(user)
}

// GetUsers retrieves all users in the database.
//
// [GET /users]
func GetUsers(event models.MyEvent, claims models.MyClaims) models.MyResponse {
	// Params

	// Connection
	db, err := helpers.Connection()
	if err != nil {
		return models.SomethingWentWrongResponse()
	}
	defer db.Close()

	// Authorized
	authorized := Authorize(db, claims.UserId)
	if !authorized {
		return models.UnauthorizedResponse()
	}

	// User not found
	users, err := repository.FindUsers(db)
	if err != nil {
		return models.SomethingWentWrongResponse()
	}

	return models.UsersFoundResponse(users)
}
