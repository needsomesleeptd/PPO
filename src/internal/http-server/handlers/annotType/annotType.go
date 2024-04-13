package annot_type_handler

import (
	service "annotater/internal/bl/anotattionTypeService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

var (
	ErrDecodingJson     = errors.New("broken request")
	ErrAddingAnnoType   = errors.New("error adding annotattion type")
	ErrGettingAnnoType  = errors.New("error getting annotattion type")
	ErrDeletingAnnoType = errors.New("error deleting annotattion type")
)

type RequestAnnotType struct {
	Description string `json:"description"`
	ClassName   string `json:"class_name"`
}

type RequestID struct {
	ID uint64 `json:"id"`
}

type ResponseModifuByID struct {
	response.Response
	models_dto.MarkupType
}

func AddAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAnnotType
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}
		markupType := models.MarkupType{
			CreatorID:   int(userID),
			Description: req.Description,
			ClassName:   req.ClassName,
		}
		err = annoTypeSevice.AddAnottationType(&markupType)
		if err != nil {
			render.JSON(w, r, response.Error(ErrAddingAnnoType.Error()))
		}
		render.JSON(w, r, response.OK())
	}
}

func GetAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error())) //TODO:: add logging here
			return
		}
		var markUp *models.MarkupType
		markUp, err = annoTypeSevice.GetAnottationTypeByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(ErrAddingAnnoType.Error()))
		}
		resp := ResponseModifuByID{MarkupType: *models_dto.ToDaMarkupType(*markUp), Response: response.OK()}
		render.JSON(w, r, resp)
	}
}

func DeleteAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingJson.Error()))
			return
		}
		err = annoTypeSevice.DeleteAnotattionType(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDeletingAnnoType.Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}
