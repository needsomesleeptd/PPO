package repository

import "annotater/internal/models"

type IAnotattionTypeRepository interface {
	AddAnottationType(markUp *models.MarkupType) error
	DeleteAnotattionType(id uint64) error
	GetAnottationTypeByID(id uint64) (*models.MarkupType, error)
	GetAnottationTypesByIDs(id []uint64) ([]models.MarkupType, error)
	GetAnottationTypesByUserID(creator_id uint64) ([]models.MarkupType, error)
	GetAllAnottationTypes() ([]models.MarkupType, error)
}
