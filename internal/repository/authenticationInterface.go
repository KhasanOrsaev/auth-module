package repository

type AuthenticationInterface interface {
	Authenticate(args ...string) (uint,error)
	Authorize(scopes []string,args ...string)(bool, error)
	GenerateToken(login,password string) (string,error)
}

