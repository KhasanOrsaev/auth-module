package basic

import (
	"encoding/base64"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)


func TestBasic_Authenticate(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()
	login := "user"
	pass:= "pass"
	id := uuid.New()
	token := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
	secret := GenerateUserSecret(login,pass)
	rows := []string{"id","name", "password_hash", "role","created_at", "updated_at", "active"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\"").WillReturnRows(sqlmock.NewRows(rows).AddRow(id, login, secret,
		nil, time.Now(),time.Now(),true))
	gormDB, _ := gorm.Open("postgres", db)
	client := Client(gormDB.LogMode(true))
	getID, err := client.Authenticate(token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, getID, id)
}

func TestBasic_Authorize(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()
	login := "user"
	pass:= "pass"
	token := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
	scopes := pq.StringArray{"user:read"}
	userRows := []string{"role"}
	mock.ExpectQuery("SELECT role FROM \"auth_users\"").WillReturnRows(sqlmock.NewRows(userRows).
		AddRow(1))
	mock.ExpectQuery("SELECT scopes FROM \"auth_roles\"").WillReturnRows(sqlmock.NewRows([]string{"scopes"}).
		AddRow(scopes))
	gormDB, _ := gorm.Open("postgres", db)
	client := Client(gormDB.LogMode(true))
	isAuthorized, err := client.Authorize(token,scopes)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, isAuthorized, true)
}