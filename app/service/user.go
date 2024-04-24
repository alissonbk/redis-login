package service

import (
	"com.github.alissonbk/go-rest-template/app/constant"
	"com.github.alissonbk/go-rest-template/app/exception"
	"com.github.alissonbk/go-rest-template/app/model/entity"
	"com.github.alissonbk/go-rest-template/app/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) GetAll() []entity.User {
	return s.repository.FindAllUser()
}

func (s *UserService) Save(user entity.User) entity.User {
	encryptedPassword, err := s.encryptPassword(user.Password)
	if err != nil {
		exception.PanicException(constant.UnknownError, "Could not encrypt the password")
	}
	user.Password = encryptedPassword
	user.Role = "ROLE_TEST"
	savedUser := s.repository.Save(&user)
	return savedUser
}

func (s *UserService) GetByID(id int) entity.User {
	user := s.repository.FindUserById(id)
	return user
}

func (s *UserService) GetByEmail(email string) entity.User {
	user := s.repository.FindUserByEmail(email)
	return user
}

func (s *UserService) Update(user entity.User) {
	s.repository.Update(user)
}

func (s *UserService) Delete(id int) {
	s.repository.DeleteUserById(id)
}

func (s *UserService) encryptPassword(passwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(passwd, 12)
}
