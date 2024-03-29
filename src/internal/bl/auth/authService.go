package service

import (
	userRepo "annotater/internal/bl/userService/userRepo"
	"annotater/internal/models"
	auth_utils "annotater/internal/pkg/authUtils"

	"github.com/pkg/errors"
)

var (
	NO_LOGIN_ERR          = errors.New("Login cannot be empty")
	NO_PASSWD_ERR         = errors.New("Password cannot be empty")
	GETTING_USER_DATA_ERR = errors.New("There is a user with this login already")
	CREATING_USER_ERR     = errors.New("Error in creating user")
	GENERATING_TOKEN_ERR  = errors.New("Error in generating token for user")
	GENERATING_HASH_ERR   = errors.New("Error in generating passwdHash for user")
	CMP_PASSED_HASH_ERR   = errors.New("Error in comparing hash and passwd")
)

type IAuthService interface {
	SignIn(candidate *models.User) (tokenStr string, err error)
	Auth(candidate *models.User) error
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

func (serv *AuthService) Auth(candidate *models.User) error {
	var passHash string
	var err error
	if candidate.Login == "" {
		return NO_LOGIN_ERR
	}

	if candidate.Password == "" {
		return NO_PASSWD_ERR
	}

	passHash, err = serv.passwordHasher.GenerateHash(candidate.Password)
	if err != nil {
		return errors.Wrap(err, GENERATING_HASH_ERR.Error())
	}
	candidate.Password = passHash

	err = serv.userRepo.CreateUser(*candidate)
	if err != nil {
		return errors.Wrap(err, CREATING_USER_ERR.Error())
	}
	return nil
}

func (serv *AuthService) SignIn(candidate *models.User) (tokenStr string, err error) {
	var user *models.User
	if candidate.Login == "" {
		return "", NO_LOGIN_ERR
	}

	if candidate.Password == "" {
		return "", NO_PASSWD_ERR
	}
	user, err = serv.userRepo.GetUserByLogin(candidate.Login)
	if err != nil {
		return "", errors.Wrap(err, GETTING_USER_DATA_ERR.Error())
	}
	err = serv.passwordHasher.ComparePasswordhash(candidate.Password, user.Password)
	if err != nil {
		return "", errors.Wrap(err, CMP_PASSED_HASH_ERR.Error())
	}
	tokenStr, err = serv.tokenizer.GenerateToken(*candidate, serv.key)
	if err != nil {
		return "", errors.Wrap(err, GENERATING_TOKEN_ERR.Error())
	}
	return tokenStr, nil
}
