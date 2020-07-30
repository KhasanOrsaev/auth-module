package basic

import (
	"auth-module/internal/models"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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
func (c *Basic) Authenticate(tokenString string) (uuid.UUID,error) {
	// получение пароля пользователя
	tokenDecode, err := base64.StdEncoding.DecodeString(tokenString)
	if err!=nil {
		return [16]byte{}, err
	}
	token := strings.Split(string(tokenDecode), ":")
	clientID := token[0]
	clientSecret := token[1]
	userSecret:= GenerateUserSecret(clientID, clientSecret)
	user := models.User{}
	c.db.Where(&models.User{Name: clientID, PasswordHash: userSecret}).Take(&user)
	if !user.Active {
		return [16]byte{}, errors.New("user is not active")
	}
	return user.ID, nil
}
// Authorize Authorize user by token
func (c *Basic) Authorize(tokenString string, scopes []string)(bool, error) {
	// получение пароля пользователя
	tokenDecode, err := base64.StdEncoding.DecodeString(tokenString)
	if err!=nil {
		return false, err
	}
	token := strings.Split(string(tokenDecode), ":")
	clientID := token[0]
	clientSecret := token[1]
	userSecret:= GenerateUserSecret(clientID, clientSecret)
	var roleID int64
	err = c.db.Model(&models.User{}).Select("role").Where("name=? and password_hash=? and active=true",
		clientID, userSecret).Row().Scan(&roleID)
	if err != nil {
		return false, err
	}
	userScopes := pq.StringArray{}
	err = c.db.Model(&models.Role{}).Select("scopes").Where("id=?").Row().Scan(&userScopes)
	if err != nil {
		return false, err
	}
	allowed := true
	for scope := range scopes {
		isFound := false
		for i := range userScopes {
			if scope == i {
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