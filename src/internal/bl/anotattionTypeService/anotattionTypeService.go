package service

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
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

func (serv *AnotattionTypeService) AddAnottationType(anotattionType *models.MarkupType) error {
	err := serv.repo.AddAnottationType(anotattionType)
	if err != nil {
		return errors.Wrap(err, "Error in adding anotattion")
	}
	return err
}

func (serv *AnotattionTypeService) DeleteAnotattionType(id uint64) error {
	err := serv.repo.DeleteAnotattionType(id)
	if err != nil {
		return errors.Wrap(err, "Error in deleting anotattion")
	}
	return err
}

func (serv *AnotattionTypeService) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	markup, err := serv.repo.GetAnottationTypeByID(id)
	if err != nil {
		return markup, errors.Wrap(err, "Error in getting anotattion")
	}
	return markup, err
}
