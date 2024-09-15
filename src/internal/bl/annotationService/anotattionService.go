package service

import (
	repository "annotater/internal/bl/annotationService/annotattionRepo"
	"annotater/internal/models"
	"bytes"
	"image"
	"image/png"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	ADDING_ANNOT_ERR_STRF          = "error in adding anotattion svc  id %v"
	ADDING_ANNOT_ERR_CREATOR_STRF  = "error in adding anotattion svc  with creator id %v"
	GETTING_ANNOT_ERR_CREATOR_STRF = "error in getting anotattion svc  with creator id %v"
	DELETING_ANNOT_ERR_STRF        = "error in deleting anotattion svc  id %v"
	GETTING_ANNOT_ERR_STRF         = "error in getting anotattion svc id %v"
	GETTING_ALL_ANNOT_ERR_STR      = "error in getting all anotattions svc"
)

var (
	ErrBoundingBoxes   = models.NewUserErr("Invalid markups bounding boxes")
	ErrInvalidFileType = models.NewUserErr("Invalid filetype")
)

type IAnotattionService interface {
	AddAnottation(anotattion *models.Markup) error
	DeleteAnotattion(id uint64) error
	GetAnottationByID(id uint64) (*models.Markup, error)
	GetAnottationByUserID(userID uint64) ([]models.Markup, error)
	GetAllAnottations() ([]models.Markup, error)
}

type AnotattionService struct {
	logger *logrus.Logger
	repo   repository.IAnotattionRepository
}

func NewAnnotattionService(loggerSrc *logrus.Logger, pRep repository.IAnotattionRepository) IAnotattionService {
	image.RegisterFormat("png", "\x89PNG\r\n\x1a\n", png.Decode, png.DecodeConfig) //for checking file formats
	return &AnotattionService{
		repo:   pRep,
		logger: loggerSrc,
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

		serv.logger.WithFields(
			logrus.Fields{"userID": anotattion.CreatorID,
				"src": "AnnotService.Add"}).
			Info("wrong bounding box format")

		return errors.Wrapf(ErrBoundingBoxes, ADDING_ANNOT_ERR_CREATOR_STRF, anotattion.CreatorID)
	}

	err := CheckPngFile(anotattion.PageData)
	if err != nil {

		serv.logger.WithFields(
			logrus.Fields{"userID": anotattion.CreatorID,
				"src": "AnnotService.Add"}).
			Info("wrong file format")
		return errors.Wrapf(ErrInvalidFileType, ADDING_ANNOT_ERR_CREATOR_STRF, anotattion.CreatorID) //maybe user wants to get why his file is broken
	}

	err = serv.repo.AddAnottation(anotattion)
	if err != nil {
		logFields := serv.logger.WithFields(
			logrus.Fields{"userID": anotattion.CreatorID,
				"src": "AnnotService.Add"})

		logFields.Warn(err)

		return errors.Wrapf(err, ADDING_ANNOT_ERR_CREATOR_STRF, anotattion.CreatorID)
	}

	serv.logger.WithFields(
		logrus.Fields{"userID": anotattion.CreatorID,
			"src": "AnnotService.Add"}).
		Info("successfully created annot")

	return err
}

func (serv *AnotattionService) DeleteAnotattion(id uint64) error {
	err := serv.repo.DeleteAnotattion(id)
	if err != nil {

		serv.logger.WithFields(
			logrus.Fields{
				"src":     "AnnotService.Delete",
				"annotID": id}).
			Warn(err)
		return errors.
			Wrapf(err, DELETING_ANNOT_ERR_STRF, id)
	}

	serv.logger.WithFields(
		logrus.Fields{
			"src":     "AnnotService.Delete",
			"annotID": id}).
		Info("successfully deleted annot")

	return err
}

func (serv *AnotattionService) GetAnottationByID(id uint64) (*models.Markup, error) {
	markup, err := serv.repo.GetAnottationByID(id)
	if err != nil {
		serv.logger.WithFields(
			logrus.Fields{
				"src":     "AnnotService.GetByID",
				"annotID": id}).
			Warn(err)

		return markup, errors.Wrapf(err, GETTING_ANNOT_ERR_STRF, id)
	}

	serv.logger.WithFields(
		logrus.Fields{
			"src":     "AnnotService.GetByID",
			"annotID": id}).
		Info("successfully got annot by ID")

	return markup, err
}

func (serv *AnotattionService) GetAnottationByUserID(userID uint64) ([]models.Markup, error) {
	markups, err := serv.repo.GetAnottationsByUserID(userID)
	if err != nil {

		serv.logger.WithFields(
			logrus.Fields{
				"src":    "AnnotService.GetByUserID",
				"userID": userID}).
			Warn(err)

		return nil, errors.Wrapf(err, GETTING_ANNOT_ERR_CREATOR_STRF, userID)
	}

	serv.logger.WithFields(
		logrus.Fields{
			"src":    "AnnotService.GetByUserID",
			"userID": userID}).
		Info("successfully got annot by userID")

	return markups, nil
}

func (serv *AnotattionService) GetAllAnottations() ([]models.Markup, error) {
	markups, err := serv.repo.GetAllAnottations()
	if err != nil {
		serv.logger.WithFields(
			logrus.Fields{
				"src": "AnnotService.GetAll"}).
			Warn(err)

		return nil, errors.Wrap(err, GETTING_ALL_ANNOT_ERR_STR)
	}

	serv.logger.WithFields(
		logrus.Fields{
			"src": "AnnotService.GetAll"}).
		Info("successfully got all annots")
	return markups, nil
}
