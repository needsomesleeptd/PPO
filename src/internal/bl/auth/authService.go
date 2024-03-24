package service

import (
	repository "annotater/internal/bl/auth/userRepo"
)

type IAuthService interface {
	SignUp()
}

type DocumentService struct {
	repo repository.IUserRepository
}
