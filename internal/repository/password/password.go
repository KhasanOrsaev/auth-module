package password

import (
	"auth-module/internal/models"
	"auth-module/internal/repository/basic"
	"github.com/go-errors/errors"
	"gorm.io/gorm"
)

type Password struct {
	db *gorm.DB
}

var client Password

func Client(db *gorm.DB) *Password {
	client.db = db
	return &client
}

// Authenticate Authenticate user by token
func (c *Password) Authenticate(args ...string) (uint,error) {
	if len(args) < 2 {
		return 0, errors.New("args not enough")
	}
	passwordHash := basic.GenerateUserSecret(args[0], args[1])
	user := models.User{}
	c.db.Where(&models.User{Name: args[0], PasswordHash: passwordHash}).Take(&user)
	if !user.Active {
		return 0, errors.New("user is not active")
	}
	return user.ID, nil
}

// Authorize Authorize user by token
func (c *Password) Authorize(scopes []string, args ...string)(bool, error) {
	if len(args) < 2 {
		return false, errors.New("args not enough")
	}
	passwordHash := basic.GenerateUserSecret(args[0], args[1])
	user := models.User{}
	c.db.Preload("Role").Where(&models.User{Name: args[0], PasswordHash: passwordHash}).Find(&user)
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

func (c *Password) GenerateToken(login,password string) (string,error) {
	return "",nil
}
