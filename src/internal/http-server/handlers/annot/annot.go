package annot_handler

import (
	service "annotater/internal/bl/annotationService"
	response "annotater/internal/lib/api"
	"annotater/internal/middleware/auth_middleware"
	"annotater/internal/models"
	models_dto "annotater/internal/models/dto"
	auth_utils "annotater/internal/pkg/authUtils"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

var (
	ErrDecodingRequest = errors.New("broken request")
	ErrGettingFile     = errors.New("no file got")
	ErrGettingBbs      = errors.New("no bbs got")
	ErrAddingAnnot     = errors.New("error adding annotattion")
	ErrGettingAnnot    = errors.New("error getting annotattion")
	ErrDeletingAnnot   = errors.New("error deleting annotattion")
)

const (
	AnnotFileFieldName = "annotFile"
	JsonBbsFieldName   = "jsonBbs"
)

type RequestAddAnnot struct { //other data is a file
	ErrorBB    []float32 `json:"error_bb"`
	ClassLabel uint64    `json:"class_label"`
}

type RequestID struct {
	ID uint64 `json:"id"`
}

type ResponseGetAnnot struct {
	response.Response
	models_dto.Markup
}

type ResponseGetAnnots struct {
	response.Response
	Markups []models_dto.Markup
}

func NewAnnotHandler(logSrc *logrus.Logger, servSrc service.IAnotattionService) AnnotHandler {
	return AnnotHandler{
		log:          logSrc,
		annotService: servSrc,
	}
}

type AnnotHandler struct {
	log          *logrus.Logger
	annotService service.IAnotattionService
}

func (h *AnnotHandler) AddAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestAddAnnot
		var pageData []byte
		userID := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		file, _, err := r.FormFile(AnnotFileFieldName)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		pageData, err = io.ReadAll(file)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		bbsString := r.FormValue(JsonBbsFieldName)

		err = json.Unmarshal([]byte(bbsString), &req)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}
		annot := models.Markup{
			PageData:   pageData,
			ErrorBB:    req.ErrorBB,
			ClassLabel: req.ClassLabel,
			CreatorID:  userID,
		}
		err = h.annotService.AddAnottation(&annot)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Error(err)
			return
		}
		h.log.Infof("annot with class_label %v and bbs %v was successfully added", req.ClassLabel, req.ErrorBB)
		render.JSON(w, r, response.OK())
	}

}

func (h *AnnotHandler) GetAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error()))
			h.log.Warn(err)
			return
		}
		var markUp *models.Markup
		markUp, err = h.annotService.GetAnottationByID(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		resp := ResponseGetAnnot{Markup: *models_dto.ToDtoMarkup(*markUp), Response: response.OK()}
		h.log.Infof("annot with ID %v was successfully fetched", req.ID)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotHandler) GetAllAnnots() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warn("cannot get userIDfrom jwt in middleware")
			return
		}

		markUps, err := h.annotService.GetAllAnottations()
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		resp := ResponseGetAnnots{Markups: models_dto.ToDtoMarkupSlice(markUps), Response: response.OK()}
		h.log.Infof("user with userID %v successfully got all annots\n", userID)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotHandler) GetAnnotsByUserID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warnf("cannot get userID from jwt %v in middleware", auth_utils.ExtractTokenFromReq(r))
			return
		}

		markUps, err := h.annotService.GetAnottationByUserID(userID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		resp := ResponseGetAnnots{Markups: models_dto.ToDtoMarkupSlice(markUps), Response: response.OK()}
		h.log.Infof("user with userID %v successfully got all his annots\n", userID)
		render.JSON(w, r, resp)
	}
}

func (h *AnnotHandler) DeleteAnnot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(auth_middleware.UserIDContextKey).(uint64)
		if !ok {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warn("cannot get userIDfrom jwt in middleware")
			return
		}

		var req RequestID
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			render.JSON(w, r, response.Error(ErrDecodingRequest.Error())) //TODO:: add logging here
			h.log.Warn(err)
			return
		}

		err = h.annotService.DeleteAnotattion(req.ID)
		if err != nil {
			render.JSON(w, r, response.Error(models.GetUserError(err).Error()))
			h.log.Warn(err)
			return
		}
		h.log.Infof("user with userID %v successfully deleted annot\n", userID)
		render.JSON(w, r, response.OK())
	}
}
