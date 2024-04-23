package repository

import (
	"annotater/internal/models"

	"github.com/google/uuid"
)

type IDocumentRepository interface {
	AddDocument(doc *models.Document) error
	DeleteDocumentByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.Document, error)
	GetDocumentsByCreatorID(id uint64) ([]models.Document, error)
	GetDocumentCountByCreator(id uint64) (int64, error)
}
