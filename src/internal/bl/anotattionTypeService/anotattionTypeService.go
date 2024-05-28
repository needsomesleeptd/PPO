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

		serv.log.WithFields(
			logrus.Fields{
				"src":         "AnnotTypeService.AddAnottationType",
				"annotTypeID": anotattionType.ID}).
			Info("inserting with empty classname")

		return ErrInsertingEmptyClass
	}
	err := serv.repo.AddAnottationType(anotattionType) // Here might be error or warning, important
	if err != nil {

		serv.log.WithFields(
			logrus.Fields{
				"src":         "AnnotTypeService.AddAnottationType",
				"annotTypeID": anotattionType.ID}).
			Error(err)

		return errors.Wrapf(err, ADDING_ANNOTATTION_ERR_STRF, anotattionType.ID)
	}

	serv.log.WithFields(
		logrus.Fields{
			"src":         "AnnotTypeService.AddAnottationType",
			"annotTypeID": anotattionType.ID}).
		Info(err)

	return err
}

func (serv *AnotattionTypeService) DeleteAnotattionType(id uint64) error {
	err := serv.repo.DeleteAnotattionType(id)
	if err != nil {

		serv.log.WithFields(
			logrus.Fields{
				"src":         "AnnotTypeService.DeleteAnottationType",
				"annotTypeID": id}).
			Error(err)

		return errors.Wrapf(err, DELETING_ANNOTATTION_ERR_STRF, id)
	}

	serv.log.WithFields(
		logrus.Fields{
			"src":         "AnnotTypeService.DeleteAnottationType",
			"annotTypeID": id}).
		Info("successfully deleted annotType with all it's annots")
	return err
}

func (serv *AnotattionTypeService) GetAnottationTypeByID(id uint64) (*models.MarkupType, error) {
	markup, err := serv.repo.GetAnottationTypeByID(id)
	if err != nil {

		serv.log.WithFields(
			logrus.Fields{
				"src":         "AnnotTypeService.GetAnottationTypeByID",
				"annotTypeID": id}).
			Error(err)

		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, id)
	}

	serv.log.WithFields(
		logrus.Fields{
			"src":         "AnnotTypeService.GetAnottationTypeByID",
			"annotTypeID": id}).
		Info("successfully got annotType by ID")

	return markup, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByIDs(ids []uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByIDs(ids)
	if err != nil {

		serv.log.WithFields(
			logrus.Fields{
				"src":          "AnnotTypeService.GetAnottationTypeByID",
				"annotTypeIDs": ids}).
			Error(err)

		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, ids)
	}

	serv.log.WithFields(
		logrus.Fields{
			"src":          "AnnotTypeService.GetAnottationTypeByID",
			"annotTypeIDs": ids}).
		Info("successfully got annotTypes by IDs")
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAnottationTypesByUserID(id uint64) ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAnottationTypesByUserID(id)
	if err != nil {
		serv.log.WithFields(
			logrus.Fields{
				"src":    "AnnotTypeService.GetAnottationTypesByUserID",
				"userID": id}).
			Error(err)
		return nil, errors.Wrapf(err, GETTING_ANNOTATTION_STR_ERR_STRF, id)
	}

	serv.log.WithFields(
		logrus.Fields{
			"src":    "AnnotTypeService.GetAnottationTypesByUserID",
			"userID": id}).
		Info("successfully got annotTypes by userID")
	return markupTypes, err
}

func (serv *AnotattionTypeService) GetAllAnottationTypes() ([]models.MarkupType, error) {
	markupTypes, err := serv.repo.GetAllAnottationTypes()
	if err != nil {

		serv.log.WithFields(
			logrus.Fields{
				"src": "AnnotTypeService.GetAllAnottationTypes"}).
			Error(err)

		return nil, errors.Wrap(err, "error getting all annotation types")
	}

	serv.log.WithFields(
		logrus.Fields{
			"src": "AnnotTypeService.GetAllAnottationTypes"}).
		Info("successfully got all annotTypes")

	return markupTypes, err
}
