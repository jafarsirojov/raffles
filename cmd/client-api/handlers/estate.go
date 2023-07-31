package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

func (h *handler) GetEstates(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	options := h.getQueryValue(r)

	estates, err := h.clientService.GetEstates(ctx, offset, limit, options)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.client-api.handlers.GetEstates h.clientService.GetEstates bad request")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.client-api.handlers.GetEstates h.clientService.GetEstates", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estates
}

func (h *handler) GetLuxuryEstates(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	offsetStr := r.URL.Query().Get("offset")
	limitStr := r.URL.Query().Get("limit")

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	estates, err := h.clientService.GetLuxuryEstates(ctx, offset, limit)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.client-api.handlers.GetLuxuryEstates h.clientService.GetLuxuryEstates bad request")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.client-api.handlers.GetLuxuryEstates h.clientService.GetLuxuryEstates", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estates
}

func (h *handler) GetEstateByID(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	estate, err := h.clientService.GetEstateByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.client-api.handlers.GetEstateByID h.clientService.GetEstateByID not found")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.client-api.handlers.GetEstateByID h.clientService.GetEstateByID", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = estate
}

func (h *handler) GetSearchOptions(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	var ctx = r.Context()

	options, err := h.clientService.GetSearchOptions(ctx)
	if err != nil {
		if err == errors.ErrNotFound {
			h.logger.Info("cmd.client-api.handlers.GetSearchOptions h.clientService.GetSearchOptions not found")
			response = responses.NotFound
			return
		}

		h.logger.Error("cmd.client-api.handlers.GetSearchOptions h.clientService.GetSearchOptions", zap.Error(err))
		response = responses.InternalErr
		return
	}

	response = responses.Success
	response.Payload = options
}

func (h *handler) ClearCache(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	h.clientService.ClearCache()

	response = responses.Success
}

func (h *handler) GetImageBaseURL(w http.ResponseWriter, r *http.Request) {
	var response structs.Response
	defer reply.Json(w, http.StatusOK, &response)

	baseURL := h.clientService.GetImageBaseURL()

	response = responses.Success
	response.Payload = struct{ URL string }{URL: baseURL}
}
