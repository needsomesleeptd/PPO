package repository

import "annotater/internal/models"

type IUserRepository interface {
	GetUserByLogin(login string) (*models.User, error)
	GetUserByCookie(cookie models.Cookie) (*models.User, error)
	GetUserByID(id uint64) (*models.User, error)
	UpdateUserByLogin(login string) error
	DeleteUserByLogin(login string) error
	AddUser(user models.User) error
}
