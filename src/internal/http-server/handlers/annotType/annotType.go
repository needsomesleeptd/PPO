package annot_type_handler

import (
	service "annotater/internal/bl/anotattionTypeService"
	response "annotater/internal/lib/api"
	logger_setup "annotater/internal/logger"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	auth_utils "annotater/internal/pkg/authUtils"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
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

func NewAnnotTypehandler(logSrc *logrus.Logger, servSrc service.IAnotattionTypeService) AnnotTypeHandler {
	return AnnotTypeHandler{
		log:           logSrc,
		annotTypeServ: servSrc,
	}
}

type AnnotTypeHandler struct {
	annotTypeServ service.IAnotattionTypeService
	log           *logrus.Logger
}

func (h *AnnotTypeHandler) AddAnnotType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAnnotType
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			h.log.Warnf("cannot get userID from jwt %v in middleware", auth_utils.ExtractTokenFromReq(r))
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf("unable to decode request %v:%v", req.ID, err)
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		markupType := models.MarkupType{
			CreatorID:   int(userID),
			Description: req.Description,
			ClassName:   req.ClassName,
			ID:          req.ID,
		}
		err = h.annotTypeServ.AddAnottationType(&markupType)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Warn(err.Error())
			return
		}
		h.log.Infof("user with userID %v successfully added annotType with ID %v\n", userID, req.ID)
		render.JSON(w, r, response.OK())
	}
}

func (h *AnnotTypeHandler) GetAnnotType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err.Error())
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		var markUp *models.MarkupType
		markUp, err = h.annotTypeServ.GetAnottationTypeByID(req.ID)
		if err != nil {
			h.log.Warnf("unable to get annnot type for ID %v  %v", req.ID, err.Error())
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetByID{MarkupType: *models_dto.ToDtoMarkupType(*markUp), Response: response.OK()}
		h.log.Infof("successfully got annotType with ID %v", req.ID)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotTypeHandler) GetAnnotTypesByIDs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestIDs
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err.Error())
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}
		var markUpTypes []models.MarkupType
		markUpTypes, err = h.annotTypeServ.GetAnottationTypesByIDs(req.IDs)
		if err != nil {
			h.log.Warnf("unable to get annnot type for ids %v  %v", req.IDs, err.Error())
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}
		h.log.Infof("successfully got annot types ids  %v", req.IDs)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotTypeHandler) GetAnnotTypesByCreatorID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			h.log.Warnf(logger_setup.UnableToGetUserifF, auth_utils.ExtractTokenFromReq(r))
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error())) //TODO:: add logging here
			return
		}

		markUpTypes, err := h.annotTypeServ.GetAnottationTypesByUserID(userID)
		if err != nil {
			h.log.Warnf("unable to get annotTypes by UserID %v", userID)
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}
		h.log.Infof("successfully got annot types  by creator id  %v:%v\n", userID, resp.MarkupTypes)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotTypeHandler) DeleteAnnotType() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			h.log.Warnf(logger_setup.UnableToDecodeUserReqF, err.Error())
			render.JSON(w, r, response.Error(ErrBrokenRequest.Error()))
			return
		}
		err = h.annotTypeServ.DeleteAnotattionType(req.ID)
		if err != nil {
			h.log.Warnf("unable to delete annot type %v:%v\n", req.ID, err.Error())
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		h.log.Infof("successfully deleted annot type %v\n", req.ID)
		render.JSON(w, r, response.OK())
	}
}

func (h *AnnotTypeHandler) GetAllAnnotTypes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		markUpTypes, err := h.annotTypeServ.GetAllAnottationTypes()
		if err != nil {
			h.log.Warnf("unable to get all annot types %v\n", err.Error())
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			return
		}
		resp := ResponseGetTypes{
			MarkupTypes: models_dto.ToDtoMarkupTypeSlice(markUpTypes),
			Response:    response.OK(),
		}
		h.log.Infof("successfully got annot types %v\n", resp.MarkupTypes)
		render.JSON(w, r, resp)
	}
}
