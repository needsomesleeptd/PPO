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
	ErrBrokenRequest    = errors.New("broken request")
	ErrAddingAnnoType   = errors.New("error adding annotattion type")
	ErrGettingAnnoType  = errors.New("error getting annotattion type")
	ErrDeletingAnnoType = errors.New("error deleting annotattion type")
)

type RequestAnnotType struct {
	ID          uint64 `json:"id"`
	Description string `json:"description"`
	ClassName   string `json:"class_name"`
}

type RequestID struct {
	ID uint64 `json:"id"`
}

type RequestIDs struct {
	IDs []uint64 `json:"ids"`
}

type ResponseGetByID struct {
	response.Response
	models_dto.MarkupType
}

type ResponseGetTypes struct {
	response.Response
	MarkupTypes []models_dto.MarkupType `json:"markupTypes"`
}

func AddAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAnnotType
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		markupType := models.MarkupType{
			CreatorID:   int(userID),
			Description: req.Description,
			ClassName:   req.ClassName,
			ID:          req.ID,
		}
		err = annoTypeSevice.AddAnottationType(&markupType)
		if err != nil {

			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func GetAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		var markUp *models.MarkupType
		markUp, err = annoTypeSevice.GetAnottationTypeByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetByID{MarkupType: *models_dto.ToDtoMarkupType(*markUp), Response: response.OK()}
		render.JSON(w, r, resp)
	}
}

func GetAnnotTypesByIDs(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestIDs
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		var markUpTypes []models.MarkupType
		markUpTypes, err = annoTypeSevice.GetAnottationTypesByIDs(req.IDs)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}

		render.JSON(w, r, resp)
	}
}

func GetAnnotTypesByCreatorID(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}

		markUpTypes, err := annoTypeSevice.GetAnottationTypesByUserID(userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}

		render.JSON(w, r, resp)
	}
}

func DeleteAnnotType(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}
		err = annoTypeSevice.DeleteAnotattionType(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		render.JSON(w, r, response.OK())
	}
}

func GetAllAnnotTypes(annoTypeSevice service.IAnotattionTypeService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		markUpTypes, err := annoTypeSevice.GetAllAnottationTypes()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}
		render.JSON(w, r, resp)
	}
}
