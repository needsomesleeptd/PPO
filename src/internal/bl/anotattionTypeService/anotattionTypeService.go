package service

import (
	repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	"annotater/internal/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ADDING_ANNOTATTION_ERR_STRF      = "Error in adding anotattion svc with ID %v\n"
	DELETING_ANNOTATTION_ERR_STRF    = "Error in deleting anotattion scv with ID %v\n"
	GETTING_ANNOTATTION_STR_ERR_STRF = "Error in getting anotattion svc with ID %v\n"
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
	log  *logrus.Logger
	repo repository.IAnotattionTypeRepository
}

func NewAnotattionTypeService(logSrc *logrus.Logger, pRep repository.IAnotattionTypeRepository) IAnotattionTypeService {
	return &AnotattionTypeService{
		log:  logSrc,
		repo: pRep,
	}
}

func (serv *AnotattionTypeService) AddAnottationType(anotattionType *models.MarkupType) error {
	if len(anotattionType.ClassName) == 0 {
		serv.log.Infof("error inserting empty class ID %v:%v\n", anotattionType.ID, ErrInsertingEmptyClass)
		return ErrInsertingEmptyClass
	}
	err := serv.repo.AddAnottationType(anotattionType)
	if err != nil {
		return errors.Wrapf(err, ADDING_ANNOTATTION_ERR_STRF, anotattionType.ID)
	}
	return err
}

func (serv *AnotattionTypeService) DeleteAnotattionType(id uint64) error {
	err := serv.repo.DeleteAnotattionType(id)
	if err != nil {
		return errors.Wrapf(err, DELETING_ANNOTATTION_ERR_STRF, id)
	}
	return err
}

func (serv *AnotattionTypeService) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	markup, err := serv.repo.GetAnottationTypeByID(id)
	if err != nil {
		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, id)
	}
	return markup, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByIDs(ids []uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByIDs(ids)
	if err != nil {
		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, ids)
	}
	if len(markupTypes) == 0 {
		return nil, models.ErrNotFound
	}
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByUserID(id uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByUserID(id)
	if err != nil {
		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, id)
	}
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAllAnottationTypes() ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAllAnottationTypes()
	if err != nil {
		return nil, errors.Wrap(err, "error getting all annotation types")
	}
	return markupTypes, err
}
