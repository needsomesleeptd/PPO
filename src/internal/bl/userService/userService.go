package service

import (
	repository "annotater/internal/bl/userService/userRepo"

	"github.com/pkg/errors"
)

type IUserService interface {
	ChangeUserRoleByLogin(login string) error
}

type UserService struct {
	userRepo repository.IUserRepository
}

func (serv *UserService) ChangeUserRole(login string) error {
	err := serv.userRepo.UpdateUserByLogin(login)
	if err != nil {
		return errors.Wrap(err, "Error in changing user role")
	}
	return err
}
