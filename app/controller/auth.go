package controller

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/exception"
	"com.github.alissonbk/go-rest-template/app/model/dto"
	"com.github.alissonbk/go-rest-template/app/service"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type AuthController struct {
	service *service.AuthService
}

func NewAuthController(s *service.AuthService) *AuthController {
	return &AuthController{service: s}
}

func (ac *AuthController) Login(ctx *gin.Context) {
	defer exception.PanicHandler(ctx)

	var loginDTO dto.LoginDTO
	err := ctx.BindJSON(&loginDTO)
	if err != nil {
		log.Error(err)
		exception.PanicException(constant.ParsingFailed, "")
	}

	tokenString := ac.service.Login(loginDTO.Email, loginDTO.Password)
	ctx.JSON(http.StatusOK, tokenString)
}

func (ac *AuthController) TestAuth(ctx *gin.Context) {
	defer exception.PanicHandler(ctx)
	fmt.Println("Ctx user: ", ctx.Value("user"))
	ctx.JSON(http.StatusOK, map[string]string{"message": "ok"})
}
