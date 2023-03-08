package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("SDFGjhdsfalshdfHFdsjkdsfds121232131afasdfac")

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func GetToken(userId int) string {
	expireTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "127.0.0.1",
			Subject: "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}
	return tokenString
}

func GetUserIdFromToken(JWTToken string) int {
	_,claims,err := parseJWT(JWTToken)
	if err != nil {
		return -1
	}
	return claims.UserId
}

func parseJWT(jwtToken string) (*jwt.Token,*Claims,error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(jwtToken,claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey,nil
	})

	return token,claims,err
}