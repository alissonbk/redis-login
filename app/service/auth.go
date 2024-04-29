package service

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/exception"
	"com.github.alissonbk/go-rest-template/app/model/entity"
	"com.github.alissonbk/go-rest-template/app/repository"
	"com.github.alissonbk/go-rest-template/config"
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
	redisConfig    *config.Redis
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

func (as *AuthService) ValidateTokenWithClaims(tokenString string) (*jwt.MapClaims, error) {
	claims := &jwt.MapClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return as.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return parsedToken.Claims.(*jwt.MapClaims), nil
}

func (as *AuthService) Login(username string, passwd string) string {
	user := as.userRepository.FindUserByEmail(username)

	err := bcrypt.CompareHashAndPassword(user.Password, []byte(passwd))
	if err != nil {
		logrus.Error(err)
		exception.PanicException(constant.Unauthorized, "user credentials are incorrect")
	}

	tokenString, err := as.createToken(user.Email)
	if err != nil {
		logrus.Error(err)
		exception.PanicException(constant.UnknownError, "could not create JWT token")
	}

	err = as.storeSessionRedis(user.Email, user)
	if err != nil {
		logrus.Error(err)
		exception.PanicException(constant.UnknownError, "could not save user session")
	}

	return tokenString
}

func (as *AuthService) storeSessionRedis(email string, user entity.User) error {
	ctx := config.RedisContextGetInstance().Ctx
	client := as.redisConfig.ConnectRedis()
	hashSetIdentifier := "user-session-" + email
	for k, v := range as.userToMap(user) {
		err := client.HSet(ctx, hashSetIdentifier, k, v).Err()
		if err != nil {
			return err
		}
	}
	return nil
}

func (as *AuthService) userToMap(user entity.User) map[string]string {
	return map[string]string{"email": user.Email, "name": user.Name, "role": user.Role}
}

func (as *AuthService) createToken(username string) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Millisecond * time.Duration(as.jwtExpiration)).Unix(),
	}
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims)

	tokenString, err := token.SignedString(as.secretKey)
	if err != nil {
		logrus.Fatal(err.Error())
		return "", err
	}

	return tokenString, nil
}
