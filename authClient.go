package auth_module

import (
	"auth-module/internal/models"
	"auth-module/internal/repository"
	"auth-module/internal/repository/basic"
	"auth-module/internal/repository/jwt"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

type AuthClient struct {
	Client repository.AuthenticationInterface
	db *gorm.DB
}
// NewAuthClient create new auth client
func NewAuthClient(db *gorm.DB, args ...string) (*AuthClient,error) {
	client := AuthClient{}
	client.db = db
	if len(args)>0 {
		if token := strings.Split(args[0], " "); len(token)>1 {
			switch token[0] {
			case "Bearer":
				client.Client = jwt.Client(client.db)
			case "Basic":
				client.Client = basic.Client(client.db)
			default:
				return nil, errors.New("incorrect auth header")
			}
		} else {
			switch args[0] {
			case "Bearer":
				client.Client = jwt.Client(client.db)
			case "Basic":
				client.Client = basic.Client(client.db)
			default:
				return nil, errors.New("incorrect auth header")
			}
		}
	}
	return &client, nil
}

//Authenticate Authenticate user
func (client *AuthClient) Authenticate(authHeader string) (uuid.UUID, error) {
	token := strings.Split(authHeader, " ")
	return client.Client.Authenticate(token[1])
}

//Authorize Authorize user
func (client *AuthClient) Authorize(authHeader string, scopes []string) (bool, error) {
	token := strings.Split(authHeader, " ")
	return client.Client.Authorize(token[1], scopes)
}

// NewUser create a new user
func (client *AuthClient) ApplyUser(login, password string, roleID int) *models.User  {
	secret := basic.GenerateUserSecret(login,password)
	role := models.Role{}
	client.db.First(&role, roleID)
	user := models.User{
		Name:         login,
		PasswordHash: secret,
		Role:         role,
		Active:       true,
		//CreatedAt: time.Now(),
	}

	client.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}, {Name: "password_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "active"}),
	}).Create(&user)
	return &user
}

// NewRole create a new role
func (client *AuthClient) NewRole(roleName string, scopes []string) *models.Role  {
	role := models.Role{Name: roleName, Scopes: scopes}
	client.db.Save(&role)
	return &role
}

func (client *AuthClient) GenerateToken(login, password string) (string,error) {
	return client.Client.GenerateToken(login, password)
}