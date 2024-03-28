package service

import (
	repository "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

var CHANGE_ROLE_ERROR_STR = "Error in changing user role"

type IUserService interface {
	ChangeUserRoleByLogin(login string) error
}

type UserService struct {
	userRepo repository.IUserRepository
}

func NewUserService(repo repository.IUserRepository) IUserService {
	return &UserService{
		userRepo: repo,
	}
}

func IsRolePermitted(currRole models.Role, reqRole models.Role) bool { //Depends on the order of roles
	return currRole >= reqRole
}

func (serv *UserService) ChangeUserRoleByLogin(login string) error { // Для созданяи админа, должна быть маграция бд на старте приложения
	err := serv.userRepo.UpdateUserByLogin(login)
	if err != nil {
		return errors.Wrap(err, "Error in changing user role")
	}
	return err
}
