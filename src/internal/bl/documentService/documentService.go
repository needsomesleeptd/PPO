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
	LOADING_DOCUMENT_ERR_STRF     = "Error in svc - loading document with id %v"
	CHECKING_DOCUMENT_ERR_STRF    = "Error in svc - checking document with id %v"
	REPORT_ERR_STRF               = "Error in svc - creating report with id %v"
	DOCUMENT_META_SAVE_ERR_STRF   = "Error in svc - saving metadata of a document with id %v"
	DOCUMENT_SAVE_ERR_STRF        = "Error in svc - saving document file with id %v"
	DOCUMENT_GET_ERR_STRF         = "Error in svc - getting document file with id %v"
	REPORT_GET_ERR_STRF           = "Error in svc - getting report file with id %v"
	DOCUMENT_GET_ERR_CREATOR_STRF = "Error in svc - getting document file by creator id %v"
	DOCUMENT_COUNT_ERR_STRF       = "Error in svc - getting document count by creator id %v"
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
		err := errors.Wrapf(ErrDocumentFormat, "document with name %v", documentMetaData.DocumentName)
		serv.logger.Info(err)
		return nil, err
	}
	err := serv.docRepo.AddDocument(&document)
	if err != nil {
		err = errors.Wrapf(err, DOCUMENT_SAVE_ERR_STRF, documentMetaData.ID)
		serv.logger.Warn(err)
		return nil, err
	}

	serv.logger.Infof("document file with id %s was saved\n", documentMetaData.ID)

	err = serv.docMetaRepo.AddDocument(&documentMetaData)
	if err != nil {
		err = errors.Wrapf(err, DOCUMENT_META_SAVE_ERR_STRF, documentMetaData.ID)
		serv.logger.Error(err)
		return nil, err
	}
	serv.logger.Infof("document metadata with id %s was saved\n", documentMetaData.ID)

	var errReport *models.ErrorReport
	errReport, err = serv.reportService.CreateReport(document)

	if err != nil {
		err = errors.Wrapf(err, REPORT_ERR_STRF, documentMetaData.ID)
		serv.logger.Error(err)
		return nil, err
	}
	serv.logger.Infof("report for document %s was created\n", documentMetaData.ID)

	err = serv.reportRepo.AddReport(errReport)
	if err != nil {
		err = errors.Wrapf(err, REPORT_ERR_STRF, documentMetaData.ID)
		serv.logger.Error(err)
		return nil, errors.Wrap(err, DOCUMENT_META_SAVE_ERR_STRF)
	}
	serv.logger.Infof("report for document %s was saved\n", documentMetaData.ID)
	return errReport, nil
}

func (serv *DocumentService) GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error) {
	documents, err := serv.docMetaRepo.GetDocumentsByCreatorID(creatorID)

	if err != nil {
		err = errors.Wrapf(err, DOCUMENT_GET_ERR_CREATOR_STRF, creatorID)
		serv.logger.Error(err)
		return nil, err
	}
	serv.logger.Infof("successfuly got document metadata by creator id %v\n", creatorID)

	return documents, err
}
func (serv *DocumentService) GetDocumentByID(ID uuid.UUID) (*models.DocumentData, error) {
	document, err := serv.docRepo.GetDocumentByID(ID)

	if err == models.ErrNotFound {
		err = errors.Wrapf(models.ErrNotFound, DOCUMENT_GET_ERR_STRF, ID)
		serv.logger.Error(err)
		return nil, err
	}

	if err != nil {
		err = errors.Wrapf(err, DOCUMENT_GET_ERR_STRF, ID)
		serv.logger.Error(err)
		return nil, err
	}

	serv.logger.Infof("successfuly got document by id %v\n", ID.String())

	return document, nil
}

func (serv *DocumentService) GetReportByID(ID uuid.UUID) (*models.ErrorReport, error) {
	report, err := serv.reportRepo.GetDocumentByID(ID)
	if err == models.ErrNotFound {
		err = errors.Wrapf(ErrNotFoundReport, REPORT_GET_ERR_STRF, ID)
		serv.logger.Error(err)
		return nil, err
	}
	if err != nil {
		err = errors.Wrapf(err, REPORT_GET_ERR_STRF, ID)
		serv.logger.Error(err)
		return nil, err
	}
	serv.logger.Infof("successfuly got report by id %v\n", ID.String())
	return report, nil
}

func (serv *DocumentService) GetDocumentCountByCreatorID(creatorID uint64) (int64, error) {
	count, err := serv.docMetaRepo.GetDocumentCountByCreator(creatorID)
	if err != nil {
		err = errors.Wrapf(err, DOCUMENT_COUNT_ERR_STRF, creatorID)
		serv.logger.Error(err)
		return -1, err
	}
	serv.logger.Infof("successfuly got document count by creatorID %v", creatorID)
	return count, err
}
