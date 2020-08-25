package auth_module

import (
	"database/sql/driver"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestNewAuthClient(t *testing.T) {
	_,err := NewAuthClient(nil, "Bearer 123")
	if err != nil {
		t.Error(err)
	}
	_,err = NewAuthClient(nil, "as 123")
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
	client, err := NewAuthClient(gormDB)
	if err != nil {
		t.Error(err)
	}
	mock.ExpectQuery("SELECT \\* FROM \"auth_roles\" WHERE").WillReturnRows(sqlmock.NewRows([]string{
		"id", "name", "scopes"}).AddRow(1, "test",pq.Array([]string{"event:read"})))
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO \"auth_users\"").WillReturnResult(driver.RowsAffected(0))
	user := client.ApplyUser("test", "test", 1)
	//mock.ExpectCommit()
	//mock.ExpectClose()
	fmt.Println(user)
}

func TestAuthClient_NewRole(t *testing.T) {

}