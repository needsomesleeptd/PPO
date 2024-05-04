package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	ErrDecodingJson     = errors.New("broken load document request")
	ErrLoadingDocument  = errors.New("error loading document")
	ErrCheckingDocument = errors.New("error checking document")
	ErrGettingID        = errors.New("error generating unqiue document ID")
	ErrInvalidFile      = errors.New("got invalid file")
	ErrGettingFile      = errors.New("error retriving file")
	ErrCreatingReport   = errors.New("error creating document report")
)

const (
	FILE_HEADER_KEY = "file"
)

type IDocumentHandler interface {
	LoadDocument(documentService service.IDocumentService) http.HandlerFunc
	CheckDocument(documentService service.IDocumentService) http.HandlerFunc
}

type Documenthandler struct {
	logger     *slog.Logger
	docService service.IDocumentService
}

type RequestLoadDocument struct {
	Document []byte `json:"document_data"`
}
type RequestCheckDocument struct {
	Document []byte `json:"document_data"`
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

func NewDocumentHandler(logSrc *slog.Logger, serv service.IDocumentService) Documenthandler {
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

func (h *Documenthandler) GetDocumentsMetaData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

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

		documentMetaData := models.DocumentMetaData{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentName: handler.Filename,
			CreationTime: time.Now(),
		}
		documentData := models.DocumentData{
			DocumentBytes: fileBytes,
			ID:            documentMetaData.ID,
		}

		var report *models.ErrorReport
		report, err = h.docService.LoadDocument(documentMetaData, documentData)
		if err != nil {
			render.JSON(w, r, response.Error(ErrCreatingReport.Error()))
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
