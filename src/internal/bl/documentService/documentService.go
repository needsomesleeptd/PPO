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
	"github.com/telkomdev/go-filesig"
)

const (
	LOADING_DOCUMENT_ERR_STR   = "Error in loading document"
	CHECKING_DOCUMENT_ERR_STR  = "Error in Checking document"
	REPORT_ERR_STR             = "Error in creating report"
	DOCUMENT_META_SAVE_ERR_STR = "Error in saving metadata of a document"
	DOCUMENT_SAVE_ERR_STR      = "Error in saving document file"
	DOCUMENT_GET_ERR_STR       = "Error in getting document file"
	REPORT_GET_ERR_STR         = "Error in getting report file"
)

var (
	ErrDocumentFormat   = models.NewUserErr("error document loaded in wrong format")
	ErrNotFoundDocument = models.NewUserErr("error not found document")
	ErrNotFoundReport   = models.NewUserErr("error not found report")
)

type IDocumentService interface {
	GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error)
	GetDocumentCountByCreatorID(creatorID uint64) (int64, error)
	LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error)
	GetDocumentByID(ID uuid.UUID) (*models.DocumentData, error)
	GetReportByID(ID uuid.UUID) (*models.ErrorReport, error)
}

type DocumentService struct {
	docMetaRepo   doc_repository.IDocumentMetaDataRepository
	reportRepo    rep_data_repo.IReportDataRepository
	docRepo       doc_data_repo.IDocumentDataRepository
	reportService rep_creator_service.IReportCreatorService
}

func NewDocumentService(docMetaRepoSrc doc_repository.IDocumentMetaDataRepository, docRepoSrc doc_data_repo.IDocumentDataRepository, reportRepoSrc rep_data_repo.IReportDataRepository, reportServSrc rep_creator_service.IReportCreatorService) IDocumentService {
	return &DocumentService{
		docMetaRepo:   docMetaRepoSrc,
		reportRepo:    reportRepoSrc,
		docRepo:       docRepoSrc,
		reportService: reportServSrc,
	}
}

func (serv *DocumentService) LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentBytes))
	if !isValid {
		return nil, ErrDocumentFormat
	}
	err := serv.docRepo.AddDocument(&document)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_SAVE_ERR_STR)
	}

	err = serv.docMetaRepo.AddDocument(&documentMetaData)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_META_SAVE_ERR_STR)
	}
	var errReport *models.ErrorReport
	errReport, err = serv.reportService.CreateReport(document)

	if err != nil {
		return nil, errors.Wrap(err, REPORT_ERR_STR)
	}

	err = serv.reportRepo.AddReport(errReport)
	if err != nil {
		return nil, errors.Wrap(err, DOCUMENT_META_SAVE_ERR_STR)
	}
	return errReport, nil
}

func (serv *DocumentService) GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error) {
	documents, err := serv.docMetaRepo.GetDocumentsByCreatorID(creatorID)

	if err != nil {
		return nil, err
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
	return report, nil
}

func (serv *DocumentService) GetDocumentCountByCreatorID(creatorID uint64) (int64, error) {
	count, err := serv.docMetaRepo.GetDocumentCountByCreator(creatorID)
	if err != nil {
		return -1, err
	}
	return count, err
}
