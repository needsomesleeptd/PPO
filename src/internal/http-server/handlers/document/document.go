package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	"errors"
	"io"
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
)

const (
	FILE_HEADER_KEY = "file"
)

type RequestLoadDocument struct {
	Document []byte `json:"document_data"`
}
type RequestCheckDocument struct {
	Document []byte `json:"document_data"`
}

type ResponseCheckDoucment struct {
	Response response.Response
	Markups  []models.Markup `json:"markups"`
}

func ExtractfileBytesHelper(file multipart.File) ([]byte, error) {

	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		return nil, ErrInvalidFile
	}

	return fileBytes, nil

}

func LoadDocument(documentService service.IDocumentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//var req RequestCheckDocument
		//err := render.DecodeJSON(r.Body, &req)
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		documentID, err := uuid.NewRandom()
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingID.Error()))
		}

		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			return
		}
		file, _, err := r.FormFile("file")

		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
		}

		document := models.Document{
			CreatorID:    userID,
			DocumentData: fileBytes,
			ID:           documentID,
			CreationTime: time.Now(),
			ChecksCount:  1, //Check here the document repo
		}

		err = documentService.LoadDocument(document)
		if err != nil {
			render.JSON(w, r, response.Error(ErrLoadingDocument.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}

func CheckDocument(documentService service.IDocumentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
			return
		}
		file, _, err := r.FormFile("file")

		if err != nil {
			render.JSON(w, r, response.Error(ErrGettingFile.Error()))
		}

		var fileBytes []byte
		fileBytes, err = ExtractfileBytesHelper(file)

		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
		}

		document := models.Document{
			CreatorID:    userID,
			DocumentData: fileBytes,
			CreationTime: time.Now(),
			ChecksCount:  1, //Check here the document repo
		} //Note that we are not checking documentID

		var markups []models.Markup
		markups, err = documentService.CheckDocument(document)
		if err != nil {
			render.JSON(w, r, response.Error(err.Error()))
			return
		}
		res := ResponseCheckDoucment{Markups: markups, Response: response.OK()}
		render.JSON(w, r, res)
	}
}
