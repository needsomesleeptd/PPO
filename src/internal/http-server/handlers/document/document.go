package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	logger_setup "annotater/internal/logger"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	pdf_utils "annotater/internal/pkg/pdfUtils"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrDecodingJson     = errors.New("broken load document request")
	ErrLoadingDocument  = errors.New("error loading document")
	ErrCheckingDocument = errors.New("error checking document")
	ErrGettingID        = errors.New("error generating unqiue document ID")
	ErrInvalidFile      = errors.New("got invalid file")
	ErrGettingFile      = errors.New("error retriving file")
	ErrCreatingReport   = errors.New("error creating document report")
	ErrBrokenRequest    = errors.New("broken request")

	ErrSendingFile      = errors.New("error sending file")
	ErrGettingPageCount = errors.New("error getting pagecount")
)

const (
	FILE_HEADER_KEY = "file"

	ErrGetttingMetaData = "error getting metadataDocuments"
)

type IDocumentHandler interface {
	LoadDocument(documentService service.IDocumentService) http.HandlerFunc
	CheckDocument(documentService service.IDocumentService) http.HandlerFunc
}

type Documenthandler struct {
	logger     *logrus.Logger
	docService service.IDocumentService
}

type RequestLoadDocument struct {
	Document []byte `json:"document_data"`
}
type RequestCheckDocument struct {
	Document []byte `json:"document_data"`
}

type RequestID struct {
	ID uuid.UUID `json:"ID"`
}

type ResponseCheckDoucment struct {
	Response    response.Response
	Markups     []models.Markup     `json:"markups"`
	MarkupTypes []models.MarkupType `json:"markupTypes"`
}

type ResponseGettingMetaData struct {
	Response          response.Response
	DocumentsMetaData []models.DocumentMetaData `json:"documents_metadata"`
}

type ResponseGetReport struct {
	Response    response.Response
	Markups     []models.Markup     `json:"markups"`
	MarkupTypes []models.MarkupType `json:"markupTypes"`
}

func NewDocumentHandler(logSrc *logrus.Logger, serv service.IDocumentService) Documenthandler {
	return Documenthandler{
		logger:     logSrc,
		docService: serv,
	}
}

func ExtractfileBytesHelper(file multipart.File) ([]byte, error) {

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, ErrInvalidFile
	}

	return fileBytes, nil

}

func writeBytesIntoResponse(w http.ResponseWriter, data []byte) error {
	w.Header().Set("Content-Type", http.DetectContentType(data))
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(data)))
	_, err := w.Write(data)
	if err != nil {
		return errors.Join(err, ErrSendingFile)
	}
	return nil

}

func (h *Documenthandler) GetDocumentByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.logger.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}

		document, err := h.docService.GetDocumentByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		err = writeBytesIntoResponse(w, document.DocumentBytes)
		if err != nil {
			render.JSON(w, r, response.Error(ErrSendingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func (h *Documenthandler) GetReportByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.logger.Warnf(logger_setup.UnableToDecodeUserReqF, err)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}

		report, err := h.docService.GetReportByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		err = writeBytesIntoResponse(w, report.ReportData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrSendingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func (h *Documenthandler) GetDocumentsMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		documentsMetaData, err := h.docService.GetDocumentsByCreatorID(userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}
		resp := ResponseGettingMetaData{Response: response.OK(), DocumentsMetaData: documentsMetaData}
		render.JSON(w, r, resp)
	}
}

func (h *Documenthandler) CreateReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, handler, err := r.FormFile(FILE_HEADER_KEY)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			h.logger.Error(err.Error())
			return
		}

		var pagesCount int
		pagesCount, err = pdf_utils.GetPdfPageCount(fileBytes)

		if err != nil {
			h.logger.Error(errors.Join(err, ErrGettingPageCount).Error())
			pagesCount = -1
		}

		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: handler.Filename,
			CreationTime: time.Now(),
			PageCount:    pagesCount,
		}
		documentData := models.DocumentData{
			DocumentBytes: fileBytes,
			ID:            documentMetaData.ID,
		}

		var report *models.ErrorReport
		report, err = h.docService.LoadDocument(documentMetaData, documentData)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.logger.Error(err.Error())
			return
		}

		w.Header().Set("Content-Type", http.DetectContentType(report.ReportData))
		w.Header().Set("Content-Length", fmt.Sprintf("%v", len(report.ReportData)))
		_, err = w.Write(report.ReportData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrCreatingReport.Error()))
			h.logger.Error(err.Error())
			return
		}

	}
}
