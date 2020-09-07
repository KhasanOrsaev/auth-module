package basic

import (
	"auth-module/internal/models"
	"encoding/base64"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type Basic struct {
	db *gorm.DB
}

var client Basic

func Client(db *gorm.DB) *Basic {
	client.db = db
	return &client
}
// Authenticate Authenticate user by token
func (c *Basic) Authenticate(args ...string) (uint,error) {
	if len(args) < 1 {
		return 0, errors.New("empty args")
	}
	clientID,clientSecret, err := parseToken(args[0])
	if err!=nil {
		return 0, err
	}
	user := models.User{}
	c.db.Where(&models.User{Name: *clientID, PasswordHash: *clientSecret}).Take(&user)
	if !user.Active {
		return 0, errors.New("user is not active")
	}
	return user.ID, nil
}
// Authorize Authorize user by token
func (c *Basic) Authorize(scopes []string, args ...string)(bool, error) {
	if len(args) < 1 {
		return false, errors.New("empty args")
	}
	clientID,clientSecret, err := parseToken(args[0])
	if err!=nil {
		return false, err
	}
	user := models.User{}
	c.db.Preload("Role").Where(&models.User{Name: *clientID, PasswordHash: *clientSecret}).Find(&user)
	allowed := true
	for _,scope := range scopes {
		isFound := false
		for _,i := range user.Role {
			if scope == i.Scope {
				isFound = true
				break
			}
		}
		if !isFound {
			allowed = false
			break
		}
	}
	return allowed,nil
}

func (c *Basic) GenerateToken(login,password string) (string,error) {
	return base64.StdEncoding.EncodeToString([]byte(login + ":" + password)),nil
}

func parseToken(tokenString string) (login, secret *string, err error) {
	// получение пароля пользователя
	tokenDecode, err := base64.StdEncoding.DecodeString(tokenString)
	if err!=nil {
		return nil,nil, err
	}
	token := strings.Split(string(tokenDecode), ":")
	login = &token[0]
	clientSecret := token[1]
	s := GenerateUserSecret(*login, clientSecret)
	secret = &s
	return
}