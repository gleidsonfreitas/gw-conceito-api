package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gwsistemas.com.br/helpers"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int
	Name           string
	Email          string
	PasswordDigest string
	Token          string
	ApiKey         string
	Enabled        bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type MyEvent struct {
	Method string
	Params []string
	Token  string
}

type MyResponse struct {
	Status  int
	Type    string
	Code    int
	Message string
	Result  string
}

type MyClaims struct {
	UserId int
	jwt.StandardClaims
}

// 500 - SomethingWentWrongResponse
func SomethingWentWrongResponse() MyResponse {
	return MyResponse{Status: 500, Type: "api-error", Code: 0, Message: "Something went wrong"}
}

// 400 - BadRequestResponse
func BadRequestResponse() MyResponse {
	return MyResponse{Status: 400, Type: "validation-error", Code: 0, Message: "Bad request"}
}

// 401 - UnauthorizedResponse
func UnauthorizedResponse() MyResponse {
	return MyResponse{Status: 401, Type: "authentication-error", Code: 0, Message: "Unauthorized"}
}

// 402 - UserEmailAlreadyUsedResponse
func UserEmailAlreadyUsedResponse(user *User) MyResponse {
	return MyResponse{Status: 402, Type: "validation-error", Code: 0, Message: "Email already used", Result: user.Email}
}

/* 200 - Responses */

// 200 - UserFoundResponse
func UserFoundResponse(user *User) MyResponse {
	type Response struct {
		Id int
		Name string
		Email string
		PasswordDigest string
		Enabled bool
	}

	result, err := json.Marshal(&Response{user.Id, user.Name, user.Email, user.PasswordDigest, user.Enabled})
	if err != nil {
		return SomethingWentWrongResponse()
	}

	return MyResponse{Status: 200, Type: "ok", Code: 0, Message: "User found", Result: string(result)}
}

// 200 - UserCreatedResponse
func UserCreatedResponse(user *User) MyResponse {
	result, err := json.Marshal(user)
	if err != nil {
		return SomethingWentWrongResponse()
	}

	return MyResponse{Status: 200, Type: "ok", Code: 0, Message: "User created", Result: string(result)}
}

// 200 - UserAuthenticatedResponse
func UserAuthenticatedResponse(user *User) MyResponse {
	result, err := json.Marshal(user)
	if err != nil {
		return SomethingWentWrongResponse()
	}

	return MyResponse{Status: 200, Type: "ok", Code: 0, Message: "User authenticated", Result: string(result)}
}

// 200 - UsersFoundResponse
func UsersFoundResponse(users []*User) MyResponse {
	type Response struct {
		Id int
		Name string
		Email string
		Enabled bool
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	modifiedUsers := make([]*Response, 0)
	for _, user := range users {
		modifiedUsers = append(modifiedUsers, &Response{user.Id, user.Name, user.Email, user.Enabled, user.CreatedAt, user.UpdatedAt})
	}

	result, err := json.Marshal(modifiedUsers)
	if err != nil {
		return SomethingWentWrongResponse()
	}

	return MyResponse{Status: 200, Type: "ok", Code: 0, Message: "Users found", Result: string(result)}
}

/* Functions */

// GenerateToken returns a signed JWT token string
func GenerateToken(userId int) string {
	mySigningKey := []byte("OYt_KSxiFK7x_7f5GuSfzmmXwqiLcj_vFx2I3G7R")

	// Declare the expiration time of the token
	// here, we have kept it as 259200 minutes
	expirationTime := time.Now().Add(259200 * time.Minute)
	claims := MyClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "nazar.io",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, _ := token.SignedString(mySigningKey)

	return signedString
}

// ValidateToken
func ValidateToken(tokenString string) (MyClaims, error) {
	defer helpers.TimeTrack("ValidateToken", time.Now())

	myClaims := MyClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &myClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte("OYt_KSxiFK7x_7f5GuSfzmmXwqiLcj_vFx2I3G7R"), nil
	})
	if err != nil {
		fmt.Println(err)
		return MyClaims{}, err
	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		myClaims.UserId = claims.UserId
		myClaims.ExpiresAt = claims.ExpiresAt
	}

	return myClaims, nil
}

/* User */

// User.ValidatePassword
func (user *User) ValidatePassword(password string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	if err != nil {
		return false
	}
	return true
}
