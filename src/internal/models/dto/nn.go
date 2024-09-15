package models_dto

import (
	"annotater/internal/models"
)

type Markup struct {
	ID         uint64    `json:"id"`
	PageData   []byte    `json:"page_data"`
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
	CreatorID  uint64    `json:"creator_id"`
}

func FromDtoMarkup(markup *Markup) models.Markup {

	return models.Markup{
		ID:         markup.ID,
		ClassLabel: markup.ClassLabel,
		ErrorBB:    markup.ErrorBB,
		PageData:   markup.PageData,
		CreatorID:  markup.CreatorID,
	}

}

func ToDtoMarkup(markup models.Markup) *Markup {
	return &Markup{
		ID:         markup.ID,
		ClassLabel: markup.ClassLabel,
		ErrorBB:    markup.ErrorBB,
		PageData:   markup.PageData,
		CreatorID:  markup.CreatorID,
	}
}

func FromDtoMarkupSlice(markupsDto []Markup) []models.Markup {
	if markupsDto == nil {
		return nil
	}
	markups := make([]models.Markup, len(markupsDto))
	for i := range markupsDto {
		markups[i] = FromDtoMarkup(&markupsDto[i])
	}
	return markups
}

func ToDtoMarkupSlice(markups []models.Markup) []Markup {
	if markups == nil {
		return nil
	}
	markupsDTO := make([]Markup, len(markups))
	for i := range markups {
		markupsDTO[i] = *ToDtoMarkup(markups[i])
	}
	return markupsDTO
}
