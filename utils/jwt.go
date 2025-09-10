package utils

import (
	"os"
	"time"
	"github.com/golang-jwt/jwt"
	// "github.com/golang-jwt/jwt/v5"
)

var SecurityKey = []byte(os.Getenv("ScurityKey"))

func GenrateToken(id uint, role string) (string,error){
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id":id,
		"role":role,
		"exp":time.Now().Add(24*time.Hour).Unix(),
	})
	return token.SignedString(SecurityKey)
}

func ExtractSecriteKey(token *jwt.Token) (interface{},error){
	return SecurityKey,nil
}

