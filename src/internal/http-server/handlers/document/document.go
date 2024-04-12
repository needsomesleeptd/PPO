package document_handler

import (
	service "annotater/internal/bl/documentService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

var (
	ErrDecodingJson    = errors.New("broken load document request")
	ErrLoadingDocument = errors.New("error loading document")
)

type RequestLoadDocument struct {
	Document []byte `json:"document_data"`
}
type RequestCheckDocument struct {
	Document []byte `json:"document_data"`
}

type ResponseLoadDocument struct {
	Response response.Response
}

type ResponseCheckDocument struct {
	Response response.Response
	Markups  []models_dto.Markup `json:"markups"`
}

func LoadDocument(documentService service.IDocumentService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestCheckDocument
		err := render.DecodeJSON(r.Body, &req)
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		document := models.Document{
			CreatorID:    userID,
			DocumentData: req.Document,
			ID:           uuid.NewRandom().ID(),
			CreationTime: time.Now(),
			ChecksCount:  1, //Check here the document repo
		}

		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}

		err = documentService.LoadDocument(document)
		if err != nil {
			render.JSON(w, r, response.Error(ErrLoadingDocument.Error()))
			return
		}

		render.JSON(w, r, response.OK())
	}
}
