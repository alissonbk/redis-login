package middleware

import (
	"com.github.alissonbk/go-rest-template/injection"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthRequired(i *injection.Injection) gin.HandlerFunc {
	authService := i.NewAuthService()
	userService := i.NewUserService()
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
		user := userService.GetByEmail(email)
		c.Set("user", user)
	}
}
