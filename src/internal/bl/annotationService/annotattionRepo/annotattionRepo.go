package repository

import "annotater/internal/models"

type IAnotattionRepository interface {
	AddAnottation(doc *models.Markup) error
	DeleteAnotattion(id uint64) error
	GetAnottationByID(id uint64) (*models.Markup, error)
}
