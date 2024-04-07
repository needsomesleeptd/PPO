package service

import (
	userRepo "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"
	auth_utils "annotater/internal/pkg/authUtils"

	"errors"
)

var (
	ErrNoLogin         = errors.New("Login cannot be empty")
	ErrNoPasswd        = errors.New("Password cannot be empty")
	ErrLoginOccupied   = errors.New("There is a user with this login already")
	ErrCreatingUser    = errors.New("Error in creating user")
	ErrGeneratingToken = errors.New("Error in generating token for user")
	ErrGeneratingHash  = errors.New("Error in generating passwdHash for user")
	ErrHashPasswdMatch = errors.New("Error in comparing hash and passwd")
)

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
		return "", errors.Join(ErrLoginOccupied, err)
	}
	err = serv.passwordHasher.ComparePasswordhash(candidate.Password, user.Password)
	if err != nil {
		return "", errors.Join(ErrHashPasswdMatch, err)
	}
	tokenStr, err = serv.tokenizer.GenerateToken(*candidate, serv.key)
	if err != nil {
		return "", errors.Join(ErrGeneratingToken, err)
	}
	return tokenStr, nil
}
