package repository

import (
	"annotater/internal/models"

	"github.com/google/uuid"
)

type IDocumentMetaDataRepository interface {
	AddDocument(doc *models.DocumentMetaData) error
	DeleteDocumentByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.DocumentMetaData, error)
	GetDocumentsByCreatorID(id uint64) ([]models.DocumentMetaData, error)
	GetDocumentCountByCreator(id uint64) (int64, error)
}
