package service

import (
	doc_data_repo "annotater/internal/bl/documentService/documentDataRepo"
	doc_repository "annotater/internal/bl/documentService/documentMetaDataRepo"
	rep_data_repo "annotater/internal/bl/documentService/reportDataRepo"
	rep_creator_service "annotater/internal/bl/reportCreatorService"

	"annotater/internal/models"
	"bytes"

	"github.com/pkg/errors"
	"github.com/telkomdev/go-filesig"
)

const (
	LOADING_DOCUMENT_ERR_STR   = "Error in loading document"
	CHECKING_DOCUMENT_ERR_STR  = "Error in Checking document"
	DOCUMENT_FORMAT_ERR_STR    = "Error document loaded in wrong format"
	REPORT_ERR_STR             = "Error in creating report"
	DOCUMENT_META_SAVE_ERR_STR = "Error in saving metadata of a document"
	DOCUMENT_SAVE_ERR_STR      = "Error in saving document file"
)

type IDocumentService interface {
	GetDocumentsByCreatorID(creatorID uint64) ([]models.DocumentMetaData, error)
	GetDocumentCountByCreatorID(creatorID uint64) (int64, error)
	LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error)
}

type DocumentService struct {
	docMetaRepo   doc_repository.IDocumentMetaDataRepository
	reportRepo    rep_data_repo.IReportDataRepository
	docRepo       doc_data_repo.IDocumentDataRepository
	reportService rep_creator_service.IReportCreatorService
}

func NewDocumentService(docMetaRepoSrc doc_repository.IDocumentMetaDataRepository, docRepoSrc doc_data_repo.IDocumentDataRepository, reportRepoSrc rep_data_repo.IReportDataRepository, reportServSrc rep_creator_service.ReportCreatorService) IDocumentService {
	return &DocumentService{
		docMetaRepo:   docMetaRepoSrc,
		reportRepo:    reportRepoSrc,
		docRepo:       docRepoSrc,
		reportService: &reportServSrc,
	}
}

func (serv *DocumentService) LoadDocument(documentMetaData models.DocumentMetaData, document models.DocumentData) (*models.ErrorReport, error) {

	isValid := filesig.IsPdf(bytes.NewReader(document.DocumentBytes))
	if !isValid {
		return nil, errors.New(DOCUMENT_FORMAT_ERR_STR)
	}
	err := serv.docRepo.AddDocument(&document)
	if err != nil {
		return nil, errors.New(DOCUMENT_SAVE_ERR_STR)
	}

	err = serv.docMetaRepo.AddDocument(&documentMetaData)
	if err != nil {
		return nil, errors.New(DOCUMENT_META_SAVE_ERR_STR)
	}
	var errReport *models.ErrorReport
	errReport, err = serv.reportService.CreateReport(document)

	if err != nil {
		return nil, errors.New(REPORT_ERR_STR)
	}

	err = serv.reportRepo.AddReport(errReport)
	if err != nil {
		return nil, errors.New(DOCUMENT_META_SAVE_ERR_STR)
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

func (serv *DocumentService) GetDocumentCountByCreatorID(creatorID uint64) (int64, error) {
	count, err := serv.docMetaRepo.GetDocumentCountByCreator(creatorID)
	if err != nil {
		return -1, err
	}
	return count, err
}
