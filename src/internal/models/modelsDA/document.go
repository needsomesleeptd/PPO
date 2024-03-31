package models_da //stands for data_acess

import (
	"annotater/internal/models"
	"time"
)

type Document struct {
	ID           uint64    `gorm:"primaryKey;column:id"`
	PageCount    int       `gorm:"column:page_count"`
	DocumentData []byte    `gorm:"column:document_data"`
	ChecksCount  int       `gorm:"column:checks_count"`
	CreatorID    uint64    `gorm:"column:creator_id"`
	CreationTime time.Time `gorm:"column:creation_time"`
}

func FromDaDocument(documentDa *Document) models.Document {
	return models.Document{
		ID:           documentDa.ID,
		PageCount:    documentDa.PageCount,
		DocumentData: documentDa.DocumentData,
		ChecksCount:  documentDa.ChecksCount,
		CreatorID:    documentDa.CreatorID,
		CreationTime: documentDa.CreationTime,
	}

}

func ToDaDocument(document models.Document) *Document {
	return &Document{
		ID:           document.ID,
		PageCount:    document.PageCount,
		DocumentData: document.DocumentData,
		ChecksCount:  document.ChecksCount,
		CreatorID:    document.CreatorID,
		CreationTime: document.CreationTime,
	}
}
