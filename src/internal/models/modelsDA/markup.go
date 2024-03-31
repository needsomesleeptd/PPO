package models_da //stands for data_acess

import "annotater/internal/models"

type Markup struct {
	ID         uint64    `gorm:"primaryKey;column:id"`
	PageData   []byte    `gorm:"column:page_data"` //png file -- the page data
	ErrorBB    []float32 `gorm:"column:error_bb`   //Bounding boxes in yolov8 format
	ClassLabel uint64    `gorm:"column:class_label"`
}

func FromDaMarkup(markupDa *Markup) models.Markup {
	return models.Markup{
		ID:         markupDa.ID,
		PageData:   markupDa.PageData,
		ErrorBB:    markupDa.ErrorBB,
		ClassLabel: markupDa.ClassLabel,
	}
}

// ToDaMarkup converts a markup Markup to a data access Markup
func ToDaMarkup(markup models.Markup) *Markup {
	return &Markup{
		ID:         markup.ID,
		PageData:   markup.PageData,
		ErrorBB:    markup.ErrorBB,
		ClassLabel: markup.ClassLabel,
	}
}
