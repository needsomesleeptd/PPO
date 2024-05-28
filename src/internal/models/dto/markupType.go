package models_dto //stands for data_acess

import "annotater/internal/models"

type MarkupType struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	CreatorID   int    `json:"creator_id"`
	ClassName   string `json:"class_name"`
}

func FromDtoMarkupType(markupTypeDa *MarkupType) models.MarkupType {
	return models.MarkupType{
		ID:          markupTypeDa.ID,
		Description: markupTypeDa.Description,
		CreatorID:   markupTypeDa.CreatorID,
		ClassName:   markupTypeDa.ClassName,
	}
}

func ToDtoMarkupType(markupType models.MarkupType) *MarkupType {
	return &MarkupType{
		ID:          markupType.ID,
		Description: markupType.Description,
		CreatorID:   markupType.CreatorID,
		ClassName:   markupType.ClassName,
	}
}

func ToDtoMarkupTypeSlice(markupTypes []models.MarkupType) []MarkupType {

	if markupTypes == nil {
		return nil
	}
	markupTypesDTO := make([]MarkupType, len(markupTypes))
	for i := range markupTypes {
		markupTypesDTO[i] = *ToDtoMarkupType(markupTypes[i])
	}
	return markupTypesDTO

}
