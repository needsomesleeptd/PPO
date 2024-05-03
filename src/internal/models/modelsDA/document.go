package models_da //stands for data_acess

import (
	"annotater/internal/models"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID `gorm:"primaryKey;column:id"`
	PageCount    int       `gorm:"column:page_count"`
	DocumentName string    `gorm:"column:document_name"`
	ChecksCount  int       `gorm:"column:checks_count"`
	CreatorID    uint64    `gorm:"column:creator_id;"`
	CreationTime time.Time `gorm:"column:creation_time"`
}

func FromDaDocument(documentDa *Document) models.DocumentMetaData {
	return models.DocumentMetaData{
		ID:           documentDa.ID,
		PageCount:    documentDa.PageCount,
		DocumentName: documentDa.DocumentName,
		CreatorID:    documentDa.CreatorID,
		CreationTime: documentDa.CreationTime,
	}

}

func FromDaDocumentSlice(documentsDa []Document) []models.DocumentMetaData {

	if documentsDa == nil {
		return nil
	}
	documents := make([]models.DocumentMetaData, len(documentsDa))

	for i, documentDA := range documentsDa {
		documents[i] = FromDaDocument(&documentDA)

	}
	return documents

}

func ToDaDocument(document models.DocumentMetaData) *Document {
	return &Document{
		ID:           document.ID,
		PageCount:    document.PageCount,
		DocumentName: document.DocumentName,
		CreatorID:    document.CreatorID,
		CreationTime: document.CreationTime,
	}
}
