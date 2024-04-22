package router

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/model/dto"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// DI
	injection := NewInjection()
	userController := injection.NewUserController()
	authController := injection.NewAuthController()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.BuildResponse[any](constant.Success, "Hello", nil))
	})
	api := router.Group("/api/v1")
	{
		// The User domain it's only for example purpose...
		user := api.Group("/user")
		user.GET("", userController.GetAll)
		user.POST("", userController.Save)
		user.GET("/:userID", userController.GetByID)
		user.PUT("/:userID", userController.Update)
		user.DELETE("/:userID", userController.Delete)

		login := api.Group("/login")
		login.POST("", authController.Login)
	}

	return router
}
