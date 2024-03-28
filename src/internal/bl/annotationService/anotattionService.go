package service

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

var ADDING_ANNOT_ERR_STR = "Error in adding anotattion"
var DELETING_ANNOT_ERR_STR = "Error in deleting anotattion"
var GETTING_ANNOT_ERR_STR = "Error in getting anotattion"

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

func (serv *AnotattionService) AddAnottation(anotattion *models.Markup) error {
	err := serv.repo.AddAnottation(anotattion)

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
