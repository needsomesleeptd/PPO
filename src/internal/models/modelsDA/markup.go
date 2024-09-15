package models_da //stands for data_acess

import (
	"annotater/internal/models"
	"encoding/json"
	"errors"

	"github.com/jackc/pgtype"
)

var (
	ErrUnMarshallMarkup = errors.New("erorr in unmarshalling markup")
	ErrMarshallMarkup   = errors.New("erorr in marshalling markup")
	ErrSettingMarkup    = errors.New("erorr in setting markup")
)

type Markup struct {
	ID         uint64       `gorm:"primaryKey;column:id"`
	PageData   []byte       `gorm:"column:page_data"`                                 //png file -- the page data
	ErrorBB    pgtype.JSONB `gorm:"type:jsonb;default:'[]';not null;column:error_bb"` //because gorm cannot store slices directly(((
	ClassLabel uint64       `gorm:"column:class_label;foreignKey:MarkupTypeID"`
	CreatorID  uint64       `gorm:"column:creator_id;foreignKey:UserID"`
}

func FromDaMarkup(markupDa *Markup) (models.Markup, error) {
	markup := models.Markup{
		ID:         markupDa.ID,
		PageData:   markupDa.PageData,
		ClassLabel: markupDa.ClassLabel,
		CreatorID:  markupDa.CreatorID,
	}
	var errorBBsJson []float32
	err := json.Unmarshal(markupDa.ErrorBB.Bytes, &errorBBsJson)
	if err != nil {
		return models.Markup{}, errors.Join(ErrUnMarshallMarkup, err)
	}
	markup.ErrorBB = errorBBsJson
	return markup, nil

}

// ToDaMarkup converts a markup Markup to a data access Markup
func ToDaMarkup(markup models.Markup) (*Markup, error) {
	markupDa := Markup{
		ID:         markup.ID,
		PageData:   markup.PageData,
		ClassLabel: markup.ClassLabel,
		CreatorID:  markup.CreatorID,
	}
	jsonB, err := json.Marshal(markup.ErrorBB)
	if err != nil {
		return nil, errors.Join(ErrMarshallMarkup, err)
	}
	err = markupDa.ErrorBB.Set(jsonB)
	if err != nil {
		return nil, errors.Join(ErrSettingMarkup, err)
	}
	return &markupDa, nil
}

func FromDaMarkupSlice(markupsDa []Markup) ([]models.Markup, error) {
	if markupsDa == nil {
		return nil, nil
	}
	markups := make([]models.Markup, len(markupsDa))
	var err error
	for i, markupDa := range markupsDa {
		markups[i], err = FromDaMarkup(&markupDa)
		if err != nil {
			return nil, err
		}
	}
	return markups, nil
}
