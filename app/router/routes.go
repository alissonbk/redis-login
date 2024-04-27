package router

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/middleware"
	"com.github.alissonbk/go-rest-template/app/model/dto"
	"com.github.alissonbk/go-rest-template/injection"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Init() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// DI
	di := injection.NewInjection()
	userController := di.NewUserController()
	authController := di.NewAuthController()

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.BuildResponse[any](constant.Success, "Hello", nil))
	})
	api := router.Group("/api/v1")
	{
		// The User domain it's only for example purpose...
		user := api.Group("/user")
		user.Use(middleware.AuthRequired(di))
		user.GET("", userController.GetAll)
		user.POST("", userController.Save)
		user.GET("/:userID", userController.GetByID)
		user.PUT("/:userID", userController.Update)
		user.DELETE("/:userID", userController.Delete)

		login := api.Group("/login")
		login.POST("", authController.Login)

		test := api.Group("/test")
		test.Use(middleware.AuthRequired(di))
		test.GET("", authController.TestAuth)
	}

	return router
}
