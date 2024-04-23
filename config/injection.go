package config

import (
	"com.github.alissonbk/go-rest-template/app/controller"
	"com.github.alissonbk/go-rest-template/app/repository"
	"com.github.alissonbk/go-rest-template/app/service"
	"gorm.io/gorm"
)

// Injection is responsible for dependency injection for each route by returning a `Controller Object` ready to be used by the router

type Injection struct {
	db *gorm.DB
}

func NewInjection() *Injection {
	return &Injection{db: ConnectDB()}
}

func (i *Injection) NewUserController() *controller.UserController {
	r := repository.NewUserRepository(i.db)
	s := service.NewUserService(r)
	return controller.NewUserController(s)
}

func (i *Injection) NewAuthController() *controller.AuthController {
	ur := repository.NewUserRepository(i.db)
	s := service.NewAuthService(ur)
	return controller.NewAuthController(s)
}

func (i *Injection) NewAuthService() *service.AuthService {
	ur := repository.NewUserRepository(i.db)
	return service.NewAuthService(ur)
}
