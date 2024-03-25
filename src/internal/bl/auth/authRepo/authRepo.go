package repository

import "annotater/internal/models"

type IAuthRepository interface {
	AddToken(cookie models.Cookie) error
}
