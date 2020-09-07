package jwt

import (
	"auth-module/internal/models"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
)

type Claims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

type JWT struct {
	secretKey []byte
	db *gorm.DB
}

var client JWT

func init() {
	v := viper.New()
	replacer := strings.NewReplacer(".","_")
	v.SetEnvKeyReplacer(replacer)

	v.SetDefault("auth.secret.key", "secret_key")
	v.BindEnv("auth.secret.key")

	client.secretKey = []byte(v.GetString("auth.secret.key"))
}

func Client(db *gorm.DB) *JWT {
	client.db = db
	return &client
}
// Authenticate Authenticate user by token
func (c *JWT) Authenticate(args ...string) (uint,error) {
	claims := &Claims{}
	if len(args) < 1 {
		return 0, errors.New("empty args")
	}
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.ParseWithClaims(args[0], claims, func(token *jwt.Token) (interface{}, error) {
		// validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
		}

		// secretKey is a []byte containing your secret
		return c.secretKey, nil
	})

	if err != nil {
		return 0, err
	}


	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	return claims.UserID, nil
}

// Authorize Authorize user by token
func (c *JWT) Authorize(scopes []string, args ...string)(bool, error) {
	if len(args) < 1 {
		return false, errors.New("empty args")
	}
	id,err := c.Authenticate(args[0])
	if err != nil {
		return false, err
	}
	user := models.User{}
	c.db.Preload("Role").Find(&user, id)
	allowed := true
	for _,scope := range scopes {
		isFound := false
		for _,i := range user.Role {
			if scope == i.Scope {
				isFound = true
				break
			}
		}
		if !isFound {
			allowed = false
			break
		}
	}
	return allowed,nil
}

func (c *JWT) GenerateToken(login,password string) (string,error) {
	user := models.User{}
	c.db.Where(&models.User{Name: login}).First(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(c.secretKey)
}