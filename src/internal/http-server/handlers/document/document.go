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

func (h *Documenthandler) LoadDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		documentID, err := uuid.NewRandom()
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingID.Error()))
			h.logger.Error(err.Error())
		}

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, _, err := r.FormFile(FILE_HEADER_KEY)

		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			h.logger.Error(err.Error())
			render.JSON(w, r, response.Error(err.Error()))
		}

		document := models.Document{
			CreatorID:    userID,
			DocumentData: fileBytes,
			ID:           documentID,
			CreationTime: time.Now(),
			ChecksCount:  1, //Check here the document repo
		}

		err = h.docService.LoadDocument(document)
		if err != nil {
			h.logger.Error(err.Error())
			render.JSON(w, r, response.Error(ErrLoadingDocument.Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func (h *Documenthandler) CheckDocument() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, _, err := r.FormFile(FILE_HEADER_KEY)

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

		var documentsCount int64
		documentsCount, err = h.docService.GetDocumentCountByCreatorID(userID)

		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}

		document := models.Document{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentData: fileBytes,
			CreationTime: time.Now(),
			ChecksCount:  int(documentsCount), //Check here the document repo
		} //Note that we are not checking documentID

		err = h.docService.LoadDocument(document)

		if err != nil {
			render.JSON(w, r, response.Error(ErrLoadingDocument.Error()))
			h.logger.Error(err.Error())
			return
		}

		var markups []models.Markup
		var markupTypes []models.MarkupType
		markups, markupTypes, err = h.docService.CheckDocument(document)
		if err != nil {
			h.logger.Error(err.Error())
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		res := ResponseCheckDoucment{Markups: markups, MarkupTypes: markupTypes, Response: response.OK()}
		h.logger.Info(fmt.Sprintf("Successfully served %d markups", len(markups)))
		render.JSON(w, r, res)
	}
}

func (h *Documenthandler) GetDocumentReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}
		file, _, err := r.FormFile(FILE_HEADER_KEY)
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
		var documentsCount int64
		documentsCount, err = h.docService.GetDocumentCountByCreatorID(userID)

		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			h.logger.Error(err.Error())
			return
		}

		document := models.Document{
			ID:           uuid.New(),
			CreatorID:    userID,
			DocumentData: fileBytes,
			CreationTime: time.Now(),
			ChecksCount:  int(documentsCount), //Check here the document repo
		} //Note that we are not checking documentID

		err = h.docService.LoadDocument(document)
		if err != nil {
			render.JSON(w, r, response.Error(ErrLoadingDocument.Error()))
			h.logger.Error(err.Error())
			return
		}

		var markups []models.Markup
		var markupTypes []models.MarkupType
		markups, markupTypes, err = h.docService.CheckDocument(document)
		if err != nil {
			h.logger.Error(err.Error())
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		report, err := h.docService.CreateReport(document.ID, markups, markupTypes)
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
