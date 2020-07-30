package console

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"testing"
)

func TestMigrate(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer db.Close()


	//mock.ExpectExec("CREATE TABLE auth_roles")

	gormDB, _ := gorm.Open("postgres", db)
	err := Migrate(gormDB.LogMode(true))
	if err != nil {
		t.Error(err)
	}
}
