package service

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
)

const (
	ADDING_ANNOTATTION_ERR_STR      = "Error in adding anotattion svc"
	DELETING_ANNOTATTION_ERR_STR    = "Error in deleting anotattion scv"
	GETTING_ANNOTATTION_STR_ERR_STR = "Error in getting anotattion svc"
)

var (
	ErrInsertingEmptyClass = models.NewUserErr("class name cannot be empty")
)

type IAnotattionTypeService interface {
	AddAnottationType(anotattion *models.MarkupType) error
	DeleteAnotattionType(id uint64) error
	GetAnottationTypeByID(id uint64) (*models.MarkupType, error)
	GetAnottationTypesByIDs(ids []uint64) ([]models.MarkupType, error)
	GetAnottationTypesByUserID(creator_id uint64) ([]models.MarkupType, error)
	GetAllAnottationTypes() ([]models.MarkupType, error)
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
	if len(anotattionType.ClassName) == 0 {
		return ErrInsertingEmptyClass
	}
	err := serv.repo.AddAnottationType(anotattionType)
	if err != nil {
		return errors.Wrap(err, ADDING_ANNOTATTION_ERR_STR)
	}
	return err
}

func (serv *AnotattionTypeService) DeleteAnotattionType(id uint64) error {
	err := serv.repo.DeleteAnotattionType(id)
	if err != nil {
		return errors.Wrap(err, DELETING_ANNOTATTION_ERR_STR)
	}
	return err
}

func (serv *AnotattionTypeService) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	markup, err := serv.repo.GetAnottationTypeByID(id)
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOTATTION_STR_ERR_STR)
	}
	return markup, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByIDs(ids []uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByIDs(ids)
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOTATTION_STR_ERR_STR)
	}
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByUserID(id uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByUserID(id)
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOTATTION_STR_ERR_STR)
	}
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAllAnottationTypes() ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAllAnottationTypes()
	if err != nil {
		return nil, errors.Wrap(err, GETTING_ANNOTATTION_STR_ERR_STR)
	}
	return markupTypes, err
}
