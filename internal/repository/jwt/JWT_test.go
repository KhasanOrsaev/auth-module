package jwt

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"testing"
	"time"
)

func init() {
	os.Setenv("AUTH_SECRET_KEY", "abc")
}

var (
	id = uint(1)
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.VWEPoFGPKV8q5vcQefQy28zhVeLUZmjSj5SGoD1VQJI"
)
func TestJWT_Authenticate(t *testing.T) {
	client := Client(nil)
	userID, err := client.Authenticate(token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, userID, id)
}

func TestJWT_Authorize(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()


	scopesWant := pq.StringArray{"event:read"}
	rows := []string{"id"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\"").WillReturnRows(sqlmock.NewRows(rows).AddRow(id))
	mock.ExpectQuery("SELECT \\* FROM \"user_permissions\" WHERE \"user_permissions\"\\.\"user_id\" = ").
		WillReturnRows(sqlmock.NewRows([]string{"role_id", "user_id"}).AddRow(1,1).AddRow(2,1))
	mock.ExpectQuery("SELECT \\* FROM \"auth_roles\" WHERE \"auth_roles\".\"id\" IN").
		WillReturnRows(sqlmock.NewRows([]string{"id","scope"}).AddRow(1,"event:read").AddRow(2,"event:write"))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
	isAuthorized, err := client.Authorize(scopesWant, token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, isAuthorized, true)
}

func TestJWT_GenerateToken(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()
	id := 1
	rows := []string{"id","name", "password_hash", "role","created_at", "updated_at", "active"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\" WHERE \"auth_users\"\\.\"name\" = \\$1").
		WillReturnRows(sqlmock.NewRows(rows).
		AddRow(id,"test", "test", nil, time.Now(),time.Now(),true))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
	newToken,err := client.GenerateToken("test","")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, newToken, token)
}