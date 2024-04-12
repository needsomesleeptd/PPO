package models_dto // stands for data_transfer_objec

import (
	"annotater/internal/models"
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID           uuid.UUID `json:"id"`
	PageCount    int       `json:"page_count"`
	DocumentData []byte    `json:"document_data"` //pdf file -- the whole file
	ChecksCount  int       `json:"checks_count"`
	CreatorID    uint64    `json:"creator_id"`
	CreationTime time.Time `json:"creation_time"`
}

func FromDtoDocument(document *Document) models.Document {

	doc := models.Document{ //TODO::Think about changing only the pointer
		ID:           document.ID,
		PageCount:    document.PageCount,
		DocumentData: document.DocumentData,
		ChecksCount:  document.ChecksCount,
		CreatorID:    document.CreatorID,
		CreationTime: document.CreationTime,
	}
	doc.DocumentData = document.DocumentData
	return doc

}

func ToDtoDocument(document models.Document) Document {
	dtoDoc := Document{
		ID:           document.ID,
		PageCount:    document.PageCount,
		DocumentData: document.DocumentData,
		ChecksCount:  document.ChecksCount,
		CreatorID:    document.CreatorID,
		CreationTime: document.CreationTime,
	}
	return dtoDoc
}
