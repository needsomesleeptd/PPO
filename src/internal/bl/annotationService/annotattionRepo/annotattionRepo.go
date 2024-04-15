package repository

import "annotater/internal/models"

type IAnotattionRepository interface {
	AddAnottation(markUp *models.Markup) error
	DeleteAnotattion(id uint64) error
	GetAnottationByID(id uint64) (*models.Markup, error)
	GetAnottationsByUserID(id uint64) ([]models.Markup, error)
	GetAllAnottations() ([]models.Markup, error)
}
