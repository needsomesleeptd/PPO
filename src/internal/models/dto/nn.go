package models_dto

import (
	"annotater/internal/models"
)

type Markup struct {
	ID         uint64    `json:"id"`
	PageData   []byte    `json:"page_data"`
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
}

func FromDtoMarkup(markup *Markup) models.Markup {

	return models.Markup{
		ID:         markup.ID,
		ClassLabel: markup.ClassLabel,
		ErrorBB:    markup.ErrorBB,
		PageData:   markup.PageData,
	}

}

func ToDtoMarkup(markup models.Markup) *Markup {
	return &Markup{
		ID:         markup.ID,
		ClassLabel: markup.ClassLabel,
		ErrorBB:    markup.ErrorBB,
		PageData:   markup.PageData,
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
