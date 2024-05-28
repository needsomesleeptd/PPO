package service

import (
	repository "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var ERROR_CHANGE_ROLE_STR = "Error in changing user role"
var ERROR_GETTING_USERS_STR = "Error in getting users"

type IUserService interface {
	ChangeUserRoleByLogin(login string, role models.Role) error
	IsRolePermitted(currRole models.Role, reqRole models.Role) bool
	GetAllUsers() ([]models.User, error)
}

type UserService struct {
	logger   *logrus.Logger
	userRepo repository.IUserRepository
}

func NewUserService(loggerSrc *logrus.Logger, repo repository.IUserRepository) IUserService {
	return &UserService{
		logger:   loggerSrc,
		userRepo: repo,
	}
}

func (serv *UserService) IsRolePermitted(currRole models.Role, reqRole models.Role) bool { //Depends on the order of roles
	return currRole >= reqRole
}

func (serv *UserService) ChangeUserRoleByLogin(login string, role models.Role) error { // Для создания админа, должна быть миграция бд на старте приложения
	user, err := serv.userRepo.GetUserByLogin(login)
	if err != nil {
		serv.logger.WithFields(
			logrus.Fields{
				"src":   "UserService.ChangeUserRoleByLogin.GetUser",
				"login": login,
				"role":  role.ToString()}).
			Error(err)
		return errors.Wrap(err, fmt.Sprintf("error changing user role with login %v wanted role %v", login, role))
	}
	user.Role = role
	err = serv.userRepo.UpdateUserByLogin(login, user)
	if err != nil {
		serv.logger.WithFields(
			logrus.Fields{
				"src":   "UserService.ChangeUserRoleByLogin.UpdateUser",
				"login": login,
				"role":  role.ToString()}).
			Error(err)
		return errors.Wrap(err, fmt.Sprintf("error changing user role with login %v wanted role %v", login, role))
	}
	serv.logger.WithFields(
		logrus.Fields{
			"src":   "UserService.ChangeUserRoleByLogin.UpdateUser",
			"login": login,
			"role":  role.ToString()}).
		Info("successfully changed userRole")
	return err
}

func (serv *UserService) GetAllUsers() ([]models.User, error) {
	users, err := serv.userRepo.GetAllUsers()
	if err != nil {
		serv.logger.WithFields(
			logrus.Fields{
				"src": "UserService.GetAllUsers"}).
			Error(err)

		return nil, errors.Wrap(err, ERROR_GETTING_USERS_STR)
	}
	serv.logger.WithFields(
		logrus.Fields{
			"src": "UserService.GetAllUsers"}).
		Info("successfully got all users data")
	return users, err
}
