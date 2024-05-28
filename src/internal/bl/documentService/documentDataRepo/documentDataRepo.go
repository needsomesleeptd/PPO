package doc_data_repo

import (
	"annotater/internal/models"

	"github.com/google/uuid"
)

type IDocumentDataRepository interface {
	AddDocument(doc *models.DocumentData) error
	DeleteDocumentByID(id uuid.UUID) error
	GetDocumentByID(id uuid.UUID) (*models.DocumentData, error)
}
