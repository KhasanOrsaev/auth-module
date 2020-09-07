package auth_module

import (
	"auth-module/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestNewAuthClient(t *testing.T) {
	_,err := NewAuthClient(nil, BearerType)
	if err != nil {
		t.Error(err)
	}
	_,err = NewAuthClient(nil, 4)
	assert.EqualError(t, err, "incorrect auth header")
}

func TestAuthClient_NewUser(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	client, err := NewAuthClient(gormDB, NullType)
	if err != nil {
		t.Error(err)
	}
	mock.ExpectQuery("INSERT INTO \"auth_users\"").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	user := client.ApplyUser("test", "test", []*models.Role{{Name: "test"},{Name: "test2"}})
	assert.Equal(t,uint(1), user.ID)
}

func TestAuthClient_NewGroup(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	client, err := NewAuthClient(gormDB, NullType)
	if err != nil {
		t.Error(err)
	}
	mock.ExpectQuery("INSERT INTO \"auth_groups\"").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	group := client.ApplyGroup("test", []*models.Role{{Name: "test"},{Name: "test2"}})
	assert.Equal(t,uint(1), group.ID)
}

func TestAuthClient_NewRole(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Error(err)
	}
	client, err := NewAuthClient(gormDB, NullType)
	if err != nil {
		t.Error(err)
	}
	roleName := "reader"
	roleScope := "event:read"
	mock.ExpectQuery("INSERT INTO \"auth_roles\"").WillReturnRows(sqlmock.NewRows(
		[]string{"id"}).AddRow(1))
	role := client.ApplyRole(roleName, roleScope)
	assert.Equal(t, uint(1), role.ID)
}