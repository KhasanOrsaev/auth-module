package jwt

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
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
	id = "8cb3e6b1-66a1-4aba-8774-f41a247f1383"
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiOGNiM2U2YjEtNjZhMS00YWJhLTg3NzQtZjQxYTI0N2YxMzgzIn0.ZHWoBaVmhVnsUiQ1r3WN98s3AmVUSQLnue_k6oIR4Kg"
)
func TestJWT_Authenticate(t *testing.T) {
	client := Client(nil)
	userID, err := client.Authenticate(token)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, userID.String(), id)
}

func TestJWT_Authorize(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()

	scopesHad := pq.StringArray{"user:read", "user:write"}
	scopesWant := pq.StringArray{"user:read"}
	rows := []string{"scopes"}
	mock.ExpectQuery("SELECT scopes FROM \"auth_roles\" JOIN auth_users ").WillReturnRows(sqlmock.NewRows(rows).
		AddRow(scopesHad))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
	isAuthorized, err := client.Authorize(token,scopesWant)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, isAuthorized, true)
}

func TestJWT_GenerateToken(t *testing.T) {
	db,mock,_ := sqlmock.New()
	defer db.Close()
	id := uuid.New()
	rows := []string{"id","name", "password_hash", "role","created_at", "updated_at", "active"}
	mock.ExpectQuery("SELECT \\* FROM \"auth_users\" WHERE \\(\"auth_users\"\\.\"name\" = \\$1\\)").
		WillReturnRows(sqlmock.NewRows(rows).
		AddRow(id,"test", "test", nil, time.Now(),time.Now(),true))
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	client := Client(gormDB)
	token,err := client.GenerateToken("test","")
	if err != nil {
		t.Error(err)
	}
	assert.NotEqual(t, "", token)
}