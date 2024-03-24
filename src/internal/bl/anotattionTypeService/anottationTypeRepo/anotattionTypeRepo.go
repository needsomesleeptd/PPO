package repository

import "annotater/internal/models"

type IAnotattionTypeRepository interface {
	AddAnottationType(doc *models.MarkupType) error
	DeleteAnotattionType(id uint64) error
	GetAnottationTypeByID(id uint64) (*models.MarkupType, error)
}
