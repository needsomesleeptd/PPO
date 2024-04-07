package repository

import "annotater/internal/models"

type IDocumentRepository interface {
	AddDocument(doc *models.Document) error
	DeleteDocumentByID(id uint64) error
	GetDocumentByID(id uint64) (*models.Document, error)
}
