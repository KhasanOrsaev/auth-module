package repository

import "github.com/google/uuid"

type AuthenticationInterface interface {
	Authenticate(tokenString string) (uuid.UUID,error)
	Authorize(tokenString string, scopes []string)(bool, error)
	GenerateToken(login,password string) (string,error)
}

