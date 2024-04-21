package models_da //stands for data_acess

import "annotater/internal/models"

type MarkupType struct {
	ID          uint64 `gorm:"primaryKey;column:id"`
	Description string `gorm:"column:description"`
	CreatorID   int    `gorm:"column:creator_id;foreignKey:UserID"`
	ClassName   string `gorm:"column:class_name"`
}

func FromDaMarkupType(markupTypeDa *MarkupType) models.MarkupType {
	return models.MarkupType{
		ID:          markupTypeDa.ID,
		Description: markupTypeDa.Description,
		CreatorID:   markupTypeDa.CreatorID,
		ClassName:   markupTypeDa.ClassName,
	}
}

// ToDaMarkupType converts a markup MarkupType to a data access MarkupType
func ToDaMarkupType(markupType models.MarkupType) *MarkupType {
	return &MarkupType{
		ID:          markupType.ID,
		Description: markupType.Description,
		CreatorID:   markupType.CreatorID,
		ClassName:   markupType.ClassName,
	}
}

func FromDaMarkupTypeSlice(markupsDa []MarkupType) []models.MarkupType {
	if markupsDa == nil {
		return nil
	}
	markups := make([]models.MarkupType, len(markupsDa))
	for i, markupDa := range markupsDa {
		markups[i] = FromDaMarkupType(&markupDa)
	}
	return markups
}
