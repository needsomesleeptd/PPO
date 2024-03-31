package models_da //stands for data_acess

import "annotater/internal/models"

type MarkupType struct {
	ID          uint64 `gorm:"primaryKey;column:id"`
	Description string `gorm:"column:description"`
	CreatorID   int    `gorm:"column:creator_id"`
}

func FromDaMarkupType(markupTypeDa *MarkupType) models.MarkupType {
	return models.MarkupType{
		ID:          markupTypeDa.ID,
		Description: markupTypeDa.Description,
		CreatorID:   markupTypeDa.CreatorID,
	}
}

// ToDaMarkupType converts a markup MarkupType to a data access MarkupType
func ToDaMarkupType(markupType models.MarkupType) *MarkupType {
	return &MarkupType{
		ID:          markupType.ID,
		Description: markupType.Description,
		CreatorID:   markupType.CreatorID,
	}
}
