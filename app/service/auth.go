package service

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/exception"
	"com.github.alissonbk/go-rest-template/app/repository"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	userRepository *repository.UserRepository
	secretKey      []byte
	jwtExpiration  int
}

func NewAuthService(ur *repository.UserRepository) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	jwtExpiration, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		logrus.Fatal(err.Error())
		return nil
	}
	return &AuthService{userRepository: ur, secretKey: []byte(secret), jwtExpiration: jwtExpiration}
}

func (as *AuthService) createToken(username string) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Millisecond * time.Duration(as.jwtExpiration)).Unix(),
		})

	tokenString, err := token.SignedString(as.secretKey)
	if err != nil {
		logrus.Fatal(err.Error())
		return "", err
	}

	return tokenString, nil
}

func (as *AuthService) validateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return as.secretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (as *AuthService) Login(username string, passwd string) string {
	user := as.userRepository.FindUserByEmail(username)
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(passwd))
	if err != nil {
		logrus.Error(err)
		exception.PanicException(constant.Unauthorized, "user credentials are incorrect")
	}

	tokenString, err := as.createToken(username)
	if err != nil {
		logrus.Error(err)
		exception.PanicException(constant.UnknownError, "could not create JWT token")
	}

	return tokenString
}
