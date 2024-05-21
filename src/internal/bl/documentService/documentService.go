package service

import (
	doc_data_repo "annotater/internal/bl/documentService/documentDataRepo"
	doc_repository "annotater/internal/bl/documentService/documentMetaDataRepo"
	rep_data_repo "annotater/internal/bl/documentService/reportDataRepo"
	rep_creator_service "annotater/internal/bl/reportCreatorService"

	"annotater/internal/models"
	"bytes"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/telkomdev/go-filesig"
)

const (
	LOADING_DOCUMENT_ERR_STR   = "Error in svc - loading document"
	CHECKING_DOCUMENT_ERR_STR  = "Error in svc - checking document"
	REPORT_ERR_STR             = "Error in svc - creating report "
	DOCUMENT_META_SAVE_ERR_STR = "Error in svc - saving metadata of a document"
	DOCUMENT_SAVE_ERR_STR      = "Error in svc - saving document file"
	DOCUMENT_GET_ERR_STR       = "Error in svc - getting document file"
	REPORT_GET_ERR_STR         = "Error in svc - getting report file"
)

var (
	ErrDocumentFormat   = models.NewUserErr("error svc document loaded in wrong format")
	ErrNotFoundDocument = models.NewUserErr("error svc not found document")
	ErrNotFoundReport   = models.NewUserErr("error svc not found report")
)

type IDocumentService interface {
	GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error)
	GetDocumentCountByCreatorID(creatorID uint64) (int64, error)
	LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error)
	GetDocumentByID(ID uuid.UUID) (*models.DocumentData, error)
	GetReportByID(ID uuid.UUID) (*models.ErrorReport, error)
}

type DocumentService struct {
	logger        *logrus.Logger
	docMetaRepo   doc_repository.IDocumentMetaDataRepository
	reportRepo    rep_data_repo.IReportDataRepository
	docRepo       doc_data_repo.IDocumentDataRepository
	reportService rep_creator_service.IReportCreatorService
}

func NewDocumentService(loggerSrc *logrus.Logger, docMetaRepoSrc doc_repository.IDocumentMetaDataRepository, docRepoSrc doc_data_repo.IDocumentDataRepository, reportRepoSrc rep_data_repo.IReportDataRepository, reportServSrc rep_creator_service.IReportCreatorService) IDocumentService {
	return &DocumentService{
		logger:        loggerSrc,
		docMetaRepo:   docMetaRepoSrc,
		reportRepo:    reportRepoSrc,
		docRepo:       docRepoSrc,
		reportService: reportServSrc,
	}
}

func (serv *DocumentService) LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentBytes))
	if !isValid {
		return nil, errors.Wrapf(ErrDocumentFormat, "document with name %v", documentMetaData.DocumentName)
	}
	err := serv.docRepo.AddDocument(&document)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_SAVE_ERR_STR)
	} else {
		serv.logger.Infof("document file with id %s was saved\n", documentMetaData.ID)
	}

	err = serv.docMetaRepo.AddDocument(&documentMetaData)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_META_SAVE_ERR_STR)
	} else {
		serv.logger.Infof("document metadata with id %s was saved\n", documentMetaData.ID)
	}
	var errReport *models.ErrorReport
	errReport, err = serv.reportService.CreateReport(document)

	if err != nil {
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	} else {
		serv.logger.Infof("report for document %s was created\n", documentMetaData.ID)
	}

	err = serv.reportRepo.AddReport(errReport)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_META_SAVE_ERR_STR)
	} else {
		serv.logger.Infof("report for document %s was saved\n", documentMetaData.ID)
	}
	return errReport, nil
}

func (serv *DocumentService) GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error) {
	documents, err := serv.docMetaRepo.GetDocumentsByCreatorID(creatorID)

	if err != nil {
		return nil, err
	} else {
		serv.logger.Infof("successfuly got document metadata by creator id %v\n", creatorID)
	}
	return documents, err
}
func (serv *DocumentService) GetDocumentByID(ID uuid.UUID) (*models.DocumentData, error) {
	document, err := serv.docRepo.GetDocumentByID(ID)

	if err == models.ErrNotFound {
		return nil, ErrNotFoundDocument
	}

	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_GET_ERR_STR)
	}

	serv.logger.Infof("successfuly got document by id %v\n", ID.String())

	return document, nil
}

func (serv *DocumentService) GetReportByID(ID uuid.UUID) (*models.ErrorReport, error) {
	report, err := serv.reportRepo.GetDocumentByID(ID)
	if err == models.ErrNotFound {
		return nil, ErrNotFoundReport
	}
	if err != nil {
		return nil, errors.Wrap(err, REPORT_GET_ERR_STR)
	}
	serv.logger.Infof("successfuly got report by id %v\n", ID.String())
	return report, nil
}

func (serv *DocumentService) GetDocumentCountByCreatorID(creatorID uint64) (int64, error) {
	count, err := serv.docMetaRepo.GetDocumentCountByCreator(creatorID)
	if err != nil {
		return -1, err
	}
	serv.logger.Infof("successfuly got document count by creatorID %v", creatorID)
	return count, err
}
