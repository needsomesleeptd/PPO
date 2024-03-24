package repository

import "annotater/internal/models"

type IDocumentRepository interface {
	AddDocument(doc *models.Document) error
	DeleteDocumentByID(uint64) error
	GetDocumentByID(doc *models.Document) (*models.Document, error)
}
