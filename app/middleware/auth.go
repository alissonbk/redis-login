package middleware

import (
	"com.github.alissonbk/go-rest-template/app/model/dto"
	"com.github.alissonbk/go-rest-template/config"
	"com.github.alissonbk/go-rest-template/injection"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

func AuthRequired(i *injection.Injection) gin.HandlerFunc {
	authService := i.NewAuthService()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		split := strings.Split(authHeader, " ")
		if len(split) != 2 {
			c.JSON(http.StatusUnauthorized, map[string]string{"message": "could not find token in the header"})
			c.Abort()
			return
		}
		token := split[1]
		if split[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, map[string]string{"message": "not a Bearer token"})
			c.Abort()
			return
		}

		claims, err := authService.ValidateTokenWithClaims(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			c.Abort()
			return
		}

		email := (*claims)["username"].(string)
		user := getUserRedis(strings.Trim(email, " "))
		c.Set("user", user)
	}
}

func getUserRedis(email string) dto.UserDTO {
	ctx := context.Background()
	redisConfig := config.Redis{}
	client := redisConfig.ConnectRedis()
	hashSetIdentifier := "user-session-" + email

	xd := client.HGetAll(ctx, hashSetIdentifier)
	fmt.Println("session from redis: ", xd)
	sessionUserMap := xd.Val()
	userDTO := dto.UserDTO{Name: sessionUserMap["name"], Email: sessionUserMap["email"], Role: sessionUserMap["role"]}
	return userDTO
}
