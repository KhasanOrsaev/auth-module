package auth_module

import (
	"auth-module/internal/models"
	"auth-module/internal/repository"
	"auth-module/internal/repository/basic"
	"auth-module/internal/repository/jwt"
	"auth-module/internal/repository/password"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	BearerType = 1
	BasicType = 2
	PasswordType = 3
	NullType = 0
)

type AuthClient struct {
	Client repository.AuthenticationInterface
	db *gorm.DB
}
// NewAuthClient create new auth client
func NewAuthClient(db *gorm.DB, clientType int) (*AuthClient,error) {
	client := AuthClient{}
	client.db = db
	switch clientType {
	case BearerType:
		client.Client = jwt.Client(client.db)
	case BasicType:
		client.Client = basic.Client(client.db)
	case PasswordType:
		client.Client = password.Client(client.db)
	case NullType:
		client.Client = nil
	default:
		return nil, errors.New("incorrect auth header")
	}


	return &client, nil
}

//Authenticate Authenticate user
func (client *AuthClient) Authenticate(args ...string) (uint, error) {
	return client.Client.Authenticate(args...)
}

//Authorize Authorize user
func (client *AuthClient) Authorize(scopes []string, args ...string) (bool, error) {
	return client.Client.Authorize(scopes, args...)
}

// ApplyUser create a new user or update
func (client *AuthClient) ApplyUser(login, password string, roles []*models.Role) *models.User  {
	secret := basic.GenerateUserSecret(login,password)
	user := models.User{
		Name:         login,
		PasswordHash: secret,
		Role:         roles,
		Active:       true,
	}
	client.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}, {Name: "password_hash"}},
		DoUpdates: clause.AssignmentColumns([]string{"role", "active"}),
	}).Omit(clause.Associations).Create(&user)
	return &user
}

// ApplyRole create a new role or update
func (client *AuthClient) ApplyRole(roleName string, scope string) *models.Role  {
	role := models.Role{Name: roleName, Scope: scope}
	client.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"scope", "updated_at"}),
	}).Create(&role)
	return &role
}

// ApplyGroup create a new group or update
func (client *AuthClient) ApplyGroup(name string, roles []*models.Role) *models.Group {
	group := models.Group{
		Name:         name,
		Role:         roles,
	}
	client.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"role"}),
	}).Omit(clause.Associations).Create(&group)
	return &group
}

func (client *AuthClient) GenerateToken(login, password string) (string,error) {
	return client.Client.GenerateToken(login, password)
}