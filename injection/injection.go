package injection

import (
	"com.github.alissonbk/go-rest-template/app/controller"
	"com.github.alissonbk/go-rest-template/app/repository"
	"com.github.alissonbk/go-rest-template/app/service"
	"com.github.alissonbk/go-rest-template/config"
	"gorm.io/gorm"
)

type Injection struct {
	db *gorm.DB
}

func NewInjection() *Injection {
	return &Injection{db: config.ConnectDB()}
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

func (i *Injection) NewUserService() *service.UserService {
	ur := repository.NewUserRepository(i.db)
	return service.NewUserService(ur)
}
