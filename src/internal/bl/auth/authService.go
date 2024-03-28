package service

import (
	authRepo "annotater/internal/bl/auth/authRepo"
	userRepo "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type IAuthService interface {
	SignIn(user *models.User) (*models.User, error)
	Auth(candidate *models.User) (*models.Cookie, error)
}

type AuthService struct {
	userRepo userRepo.IUserRepository
	authRepo authRepo.IAuthRepository
}

func (serv *AuthService) Auth(candidate models.User) (*models.User, error) {
	user, err := serv.userRepo.GetUserByLogin(candidate.Login)
	if err != nil {
		return nil, errors.Wrap(err, "Error in getting user data")
	}
	if user.Login == candidate.Login {
		return nil, errors.New("There is a user with this login already") //replace errors with const values
	}
	err = serv.userRepo.AddUser(candidate)
	if err != nil {
		return nil, errors.Wrap(err, "Error in saving user")
	}
	return &candidate, nil
}

func (serv *AuthService) SignIn(candidate *models.User) (*models.Cookie, error) {
	user, err := serv.userRepo.GetUserByLogin(candidate.Login)
	if err != nil {
		return nil, err
	}
	if candidate.Password != user.Password {
		return nil, errors.New("The passwords didn't match")
	}

	token := models.Cookie{
		UserID:  user.ID,
		ExpTime: (3600 * 60 * 48) * time.Second, // 2 days
		Token:   uuid.NewString(),
		Role:    user.Role,
	}
	err = serv.authRepo.AddToken(token)
	if err != nil {
		return nil, err
	}

	return &token, nil

}
