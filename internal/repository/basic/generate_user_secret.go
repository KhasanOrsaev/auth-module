package basic

import (
	"crypto/md5"
	"encoding/hex"
)

func GenerateUserSecret (ClientID, ClientSecret string) string{
	loginHash := md5.New()
	loginHash.Write([]byte(ClientID))
	dstLogin := make([]byte, hex.EncodedLen(len(loginHash.Sum(nil))))
	hex.Encode(dstLogin, loginHash.Sum(nil))
	PasswordHash := md5.New()
	PasswordHash.Write([]byte(ClientSecret))
	dstPassword := make([]byte, hex.EncodedLen(len(PasswordHash.Sum(nil))))
	hex.Encode(dstPassword, PasswordHash.Sum(nil))
	userSecretHash := md5.New()
	userSecretHash.Write(append(dstLogin, dstPassword...))
	return hex.EncodeToString(userSecretHash.Sum(nil))
}
