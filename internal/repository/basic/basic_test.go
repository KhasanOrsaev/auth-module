package basic

import (
	"encoding/base64"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
	"time"
)


func TestBasic_Authenticate(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()
	login := "user"
	pass:= "pass"
	id := uint(1)
	token := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
	secret := GenerateUserSecret(login,pass)
	rows := []string{"id","name", "password_hash", "role","created_at", "updated_at", "active"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\"").WillReturnRows(sqlmock.NewRows(rows).AddRow(id, login, secret,
		nil, time.Now(),time.Now(),true))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
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
	scopes := pq.StringArray{"event:read"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\" WHERE").WillReturnRows(sqlmock.NewRows([]string{"id", "active"}).
		AddRow(1, true))
	mock.ExpectQuery("SELECT \\* FROM \"user_permissions\" WHERE \"user_permissions\"\\.\"user_id\" = ").
		WillReturnRows(sqlmock.NewRows([]string{"role_id", "user_id"}).AddRow(1,1).AddRow(2,1))
	mock.ExpectQuery("SELECT \\* FROM \"auth_roles\" WHERE \"auth_roles\".\"id\" IN").
		WillReturnRows(sqlmock.NewRows([]string{"id","scope"}).AddRow(1,"event:read").AddRow(2,"event:write"))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
	isAuthorized, err := client.Authorize(scopes,token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, isAuthorized, true)
}