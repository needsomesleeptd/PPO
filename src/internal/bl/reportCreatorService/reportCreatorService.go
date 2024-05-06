package rep_creator_service

import (
	nn "annotater/internal/bl/NN"
	annot_type_repository "annotater/internal/bl/anotattionTypeService/anottationTypeRepo"
	report_creator "annotater/internal/bl/reportCreatorService/reportCreator"
	"annotater/internal/models"
	"bytes"

	"github.com/pkg/errors"
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
	annotTypeRepo annot_type_repository.IAnotattionTypeRepository
	neuralNetwork nn.INeuralNetwork
	reportWorker  report_creator.IReportCreator
}

func NewDocumentService(pNN nn.INeuralNetwork, typeRep annot_type_repository.IAnotattionTypeRepository, repCreator report_creator.IReportCreator) IReportCreatorService {
	return &ReportCreatorService{
		neuralNetwork: pNN,
		annotTypeRepo: typeRep,
		reportWorker:  repCreator,
	}
}

func (serv *ReportCreatorService) NNMarkupsReq(document models.DocumentData) ([]models.Markup, []models.MarkupType, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentBytes))
	if !isValid {
		return nil, nil, ErrDocumentFormat
	}
	markups, err := serv.neuralNetwork.Predict(document)
	if err != nil {
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	ids := make([]uint64, len(markups))
	for i := range ids {
		ids[i] = markups[i].ClassLabel
	}
	markupTypes, err := serv.annotTypeRepo.GetAnottationTypesByIDs(ids)
	if err != nil {
		return nil, nil, errors.Wrap(err, CHECKING_DOCUMENT_ERR_STR)
	}
	return markups, markupTypes, err
}

func (serv *ReportCreatorService) CreateReport(document models.DocumentData) (*models.ErrorReport, error) {

	var report *models.ErrorReport
	markups, markupTypes, err := serv.NNMarkupsReq(document)
	if err != nil {
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	}
	report, err = serv.reportWorker.CreateReport(document.ID, markups, markupTypes)
	if err != nil {
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	}
	return report, err
}
