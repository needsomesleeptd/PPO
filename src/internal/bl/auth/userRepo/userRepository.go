package repository

import "annotater/internal/models"

type IUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	UpdateUserByID(id uint64) error
	DeleteUserByID(id uint64) error
}
