package service

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
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

func (serv *AnotattionService) AddAnottation(anotattion *models.Markup) error {
	err := serv.repo.AddAnottation(anotattion)

	if err != nil {
		return errors.Wrap(err, "Error in adding anotattion")
	}
	return err
}

func (serv *AnotattionService) DeleteAnotattion(id uint64) error {
	err := serv.repo.DeleteAnotattion(id)
	if err != nil {
		return errors.Wrap(err, "Error in deleting anotattion")
	}
	return err
}

func (serv *AnotattionService) GetAnottationByID(id uint64) (*models.Markup, error) {
	markup, err := serv.repo.GetAnottationByID(id)
	if err != nil {
		return markup, errors.Wrap(err, "Error in getting anotattion")
	}
	return markup, err
}
