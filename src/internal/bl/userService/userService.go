package service

import (
	repository "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

var ERROR_CHANGE_ROLE_STR = "Error in changing user role"

type IUserService interface {
	ChangeUserRoleByLogin(login string, role models.Role) error
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

func (serv *UserService) ChangeUserRoleByLogin(login string, role models.Role) error { // Для созданяи админа, должна быть маграция бд на старте приложения
	user, err := serv.userRepo.GetUserByLogin(login)
	if err != nil {
		return errors.Wrap(err, ERROR_CHANGE_ROLE_STR)
	}
	user.Role = role
	err = serv.userRepo.UpdateUserByLogin(login, user)
	if err != nil {
		return errors.Wrap(err, ERROR_CHANGE_ROLE_STR)
	}
	return err
}
