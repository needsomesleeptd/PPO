package service

import (
	userRepo "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"
	auth_utils "annotater/internal/pkg/authUtils"

	"errors"
)

var (
	ErrNoLogin       = models.NewUserErr("login cannot be empty")
	ErrNoPasswd      = models.NewUserErr("password cannot be empty")
	ErrWrongLogin    = models.NewUserErr("wrong login")
	ErrWrongPassword = models.NewUserErr("wrong password")
)

var (
	ErrCreatingUser    = errors.New("error in creating user")
	ErrGeneratingToken = errors.New("error in generating token for user")
	ErrGeneratingHash  = errors.New("error in generating passwdHash for user")
)

const SECRET = "secret"

type IAuthService interface {
	SignIn(candidate *models.User) (tokenStr string, err error)
	SignUp(candidate *models.User) error
}

type AuthService struct {
	userRepo       userRepo.IUserRepository
	passwordHasher auth_utils.IPasswordHasher
	tokenizer      auth_utils.ITokenHandler
	key            string
}

func NewAuthService(repo userRepo.IUserRepository, hasher auth_utils.IPasswordHasher, token auth_utils.ITokenHandler, k string) IAuthService {
	return &AuthService{
		userRepo:       repo,
		passwordHasher: hasher,
		tokenizer:      token,
		key:            k,
	}
}

func (serv *AuthService) SignUp(candidate *models.User) error {
	var passHash string
	var err error
	if candidate.Login == "" {
		return ErrNoLogin
	}

	if candidate.Password == "" {
		return ErrNoPasswd
	}

	passHash, err = serv.passwordHasher.GenerateHash(candidate.Password)
	if err != nil {
		return errors.Join(ErrGeneratingHash, err)
	}
	candidateHashedPasswd := *candidate
	candidateHashedPasswd.Password = passHash

	err = serv.userRepo.CreateUser(&candidateHashedPasswd)
	if err != nil {
		return errors.Join(ErrCreatingUser, err)
	}
	return nil
}

func (serv *AuthService) SignIn(candidate *models.User) (tokenStr string, err error) {
	var user *models.User
	if candidate.Login == "" {
		return "", ErrNoLogin
	}

	if candidate.Password == "" {
		return "", ErrNoPasswd
	}
	user, err = serv.userRepo.GetUserByLogin(candidate.Login)

	if err != nil {
		return "", errors.Join(ErrWrongLogin, err)
	}
	err = serv.passwordHasher.ComparePasswordhash(candidate.Password, user.Password)
	if err != nil {
		return "", errors.Join(ErrWrongPassword, err)
	}
	tokenStr, err = serv.tokenizer.GenerateToken(*user, serv.key)
	if err != nil {
		return "", errors.Join(ErrGeneratingToken, err)
	}
	return tokenStr, nil
}
