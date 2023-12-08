package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Server) ValidateJWT(accessToken string) (token *jwt.Token, err error) {
	accessToken, err = getToken(accessToken)
	if err != nil {
		return
	}

	token, err = jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.SecretKey), nil
	})
	if err != nil {
		return
	}

	return
}

func getToken(auth string) (string, error) {
	jwtToken := strings.Split(auth, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("invalid token")
	}

	return jwtToken[1], nil
}

func (s *Server) GenerateJWT(userID int64) (token string, err error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": fmt.Sprint(userID),
	})

	token, err = claims.SignedString([]byte(s.SecretKey))
	if err != nil {
		return
	}
	return
}

func (s *Server) GetJWTClaims(token *jwt.Token, key string) string {
	return token.Claims.(jwt.MapClaims)[key].(string)
}
