package service

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

const (
	ERROR_ADDING_ANNOTATTION_STR   = "Error in adding anotattion"
	ERROR_DELETING_ANNOTATTION_STR = "Error in deleting anotattion"
	ERROR_GETTING_ANNOTATTION_STR  = "Error in deleting anotattion"
)

type IAnotattionTypeService interface {
	AddAnottationType(anotattion *models.MarkupType) error
	DeleteAnotattionType(id uint64) error
	GetAnottationTypeByID(id uint64) (*models.MarkupType, error)
}

type AnotattionTypeService struct {
	repo repository.IAnotattionTypeRepository
}

func NewAnotattionTypeService(pRep repository.IAnotattionTypeRepository) IAnotattionTypeService {
	return &AnotattionTypeService{
		repo: pRep,
	}
}

func (serv *AnotattionTypeService) AddAnottationType(anotattionType *models.MarkupType) error { //
	err := serv.repo.AddAnottationType(anotattionType)
	if err != nil {
		return errors.Wrap(err, ERROR_ADDING_ANNOTATTION_STR)
	}
	return err //create service for validation, answering if you have access or not (getting userID as an argument)
}

func (serv *AnotattionTypeService) DeleteAnotattionType(id uint64) error {
	err := serv.repo.DeleteAnotattionType(id)
	if err != nil {
		return errors.Wrap(err, ERROR_DELETING_ANNOTATTION_STR)
	}
	return err
}

func (serv *AnotattionTypeService) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	markup, err := serv.repo.GetAnottationTypeByID(id)
	if err != nil {
		return nil, errors.Wrap(err, ERROR_GETTING_ANNOTATTION_STR)
	}
	return markup, err
}
