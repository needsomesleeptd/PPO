package models_da //stands for data_acess

import (
	"annotater/internal/models"
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type BBCoordsSlice []float32 //because gorm cannot store slices directly(((

func (fs *BBCoordsSlice) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid data type for Bounding boxes")
	}
	return json.Unmarshal(bytes, &fs)
}

func (fs *BBCoordsSlice) Value() (driver.Value, error) {
	return json.Marshal(fs)
}

type Markup struct {
	ID         uint64        `gorm:"primaryKey;column:id"`
	PageData   []byte        `gorm:"column:page_data"`           //png file -- the page data
	ErrorBB    BBCoordsSlice `gorm:"type:jsonb;column:error_bb"` //Bounding boxes in yolov8 format
	ClassLabel uint64        `gorm:"column:class_label"`
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
