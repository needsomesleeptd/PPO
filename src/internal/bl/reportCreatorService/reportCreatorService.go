package rep_creator_service

import (
	nn "annotater/internal/bl/NN"
	annot_type_repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	report_creator "annotater/internal/bl/reportCreatorService/reportCreator"
	"annotater/internal/models"
	"bytes"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/telkomdev/go-filesig"
)

const (
	LOADING_DOCUMENT_ERR_STR  = "Error in loading document"
	CHECKING_DOCUMENT_ERR_STR = "Error in Checking document"
	REPORT_ERR_STR            = "Error in creating report"
)

var (
	ErrDocumentFormat = models.NewUserErr("Error document loaded in wrong format")
)

type IReportCreatorService interface {
	CreateReport(document models.DocumentData) (*models.ErrorReport, error)
}

type ReportCreatorService struct {
	logger        *logrus.Logger
	annotTypeRepo annot_type_repository.IAnotattionTypeRepository
	neuralNetwork nn.INeuralNetwork
	reportWorker  report_creator.IReportCreator
}

func NewDocumentService(loggerSrc *logrus.Logger, pNN nn.INeuralNetwork, typeRep annot_type_repository.IAnotattionTypeRepository, repCreator report_creator.IReportCreator) IReportCreatorService {
	return &ReportCreatorService{
		logger:        loggerSrc,
		neuralNetwork: pNN,
		annotTypeRepo: typeRep,
		reportWorker:  repCreator,
	}
}

func (serv *ReportCreatorService) NNMarkupsReq(document models.DocumentData) ([]models.Markup, []models.MarkupType, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentBytes))
	if !isValid {
		serv.logger.Infof("report creator svc - err in creating report with ID %v : %v", document.ID, ErrDocumentFormat)
		return nil, nil, ErrDocumentFormat
	}
	markups, err := serv.neuralNetwork.Predict(document)
	if err != nil {
		serv.logger.Errorf("report creator svc - err in getting markups for document  with ID %v : %v", document.ID, err)
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	ids := make([]uint64, len(markups))
	for i := range ids {
		ids[i] = markups[i].ClassLabel
	}
	serv.logger.Infof("report creator svc - successfully got markups for document with id %v", document.ID)
	markupTypes, err := serv.annotTypeRepo.GetAnottationTypesByIDs(ids)
	if err != nil {
		serv.logger.Warnf("report creator svc - failed to get markups for document with id %v", document.ID)
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	serv.logger.Infof("report creator svc - successfully got markupTypes for document with id %v", document.ID)
	return markups, markupTypes, err
}

func (serv *ReportCreatorService) CreateReport(document models.DocumentData) (*models.ErrorReport, error) {

	var report *models.ErrorReport
	markups, markupTypes, err := serv.NNMarkupsReq(document)
	if err != nil {
		serv.logger.Warnf("report creator svc - err in getting markups from NN : %v", err)
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	}
	report, err = serv.reportWorker.CreateReport(document.ID, markups, markupTypes)
	if err != nil {
		serv.logger.Warnf("report cretor svc - err in creating report: %v", err)
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	}
	serv.logger.Infof("report creator svc - succesfully created report  for document with ID: %v", document.ID)
	return report, err
}
