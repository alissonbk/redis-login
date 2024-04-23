package middleware

import (
	"com.github.alissonbk/go-rest-template/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthRequired(i *config.Injection) gin.HandlerFunc {
	authService := i.NewAuthService()
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		split := strings.Split(authHeader, " ")
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
		// TODO: make claims return username and user role
		fmt.Println("claims: ", claims)
	}
}
