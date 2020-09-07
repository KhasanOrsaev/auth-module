package migration

import (
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"testing"
)

func TestMigrate(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	mock.ExpectExec("CREATE TABLE \"auth_roles\"").WillReturnResult(driver.RowsAffected(0))
	mock.ExpectExec("CREATE INDEX \"idx_auth_roles_deleted_at\"").WillReturnResult(driver.RowsAffected(0))

	mock.ExpectExec("CREATE TABLE \"auth_users\"").WillReturnResult(driver.RowsAffected(0))
	mock.ExpectExec("CREATE INDEX \"user__name__password_hash__idx\"").WillReturnResult(driver.RowsAffected(0))
	mock.ExpectExec("CREATE INDEX \"idx_auth_users_deleted_at\"").WillReturnResult(driver.RowsAffected(0))

	mock.ExpectExec("CREATE TABLE \"user_permissions\"").WillReturnResult(driver.RowsAffected(0))

	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	err := Migrate(gormDB)
	if err != nil {
		t.Error(err)
	}
}
