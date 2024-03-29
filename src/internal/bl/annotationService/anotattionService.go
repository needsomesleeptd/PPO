package service

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	"bytes"
	"image"

	"github.com/pkg/errors"
)

const (
	ADDING_ANNOT_ERR_STR   = "Error in adding anotattion"
	DELETING_ANNOT_ERR_STR = "Error in deleting anotattion"
	GETTING_ANNOT_ERR_STR  = "Error in getting anotattion"
	INVALID_BBS_ERR_STR    = "Invalid markups bounding boxes"
	INVALID_FILE_ERR_STR   = "Invalid filetype" //checking only png files
)

type IAnotattionService interface {
	AddAnottation(anotattion *models.Markup) error
	DeleteAnotattion(id uint64) error
	GetAnottationByID(id uint64) (*models.Markup, error)
}

type AnotattionService struct {
	repo repository.IAnotattionRepository
}

func NewAnnotattionService(pRep repository.IAnotattionRepository) IAnotattionService {
	return &AnotattionService{
		repo: pRep,
	}
}

func AreBBsValid(slice []float32) bool { //TODO:: think if i want to export everything
	for _, value := range slice {
		if value < 0.0 || value > 1.0 {
			return false
		}
	}
	return true
}

func CheckPngFile(pngFile []byte) error {
	_, _, err := image.Decode(bytes.NewReader(pngFile))
	if err != nil {
		return err
	}
	return nil

}

func (serv *AnotattionService) AddAnottation(anotattion *models.Markup) error {
	if !AreBBsValid(anotattion.ErrorBB) {
		return errors.New(INVALID_BBS_ERR_STR)
	}

	err := CheckPngFile(anotattion.PageData)
	if err != nil {
		return errors.New(INVALID_FILE_ERR_STR)
	}

	err = serv.repo.AddAnottation(anotattion)
	if err != nil {
		return errors.Wrap(err, ADDING_ANNOT_ERR_STR)
	}
	return err
}

func (serv *AnotattionService) DeleteAnotattion(id uint64) error {
	err := serv.repo.DeleteAnotattion(id)
	if err != nil {
		return errors.Wrap(err, DELETING_ANNOT_ERR_STR)
	}
	return err
}

func (serv *AnotattionService) GetAnottationByID(id uint64) (*models.Markup, error) {
	markup, err := serv.repo.GetAnottationByID(id)
	if err != nil {
		return markup, errors.Wrap(err, GETTING_ANNOT_ERR_STR)
	}
	return markup, err
}
