package handlers

import (
	"crm/internal/responses"
	"crm/internal/structs"
	"crm/pkg/errors"
	"crm/pkg/reply"
	"encoding/json"
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

func (h *handler) getQueryValue(r *http.Request) structs.SearchOptionsDTO {
	var options structs.SearchOptionsDTO
	var err error
	priceMinStr := r.URL.Query().Get("priceMin")
	options.PriceMin, err = strconv.Atoi(priceMinStr)
	if err != nil {
		options.PriceMin = 0
	}

	priceMaxStr := r.URL.Query().Get("priceMax")
	options.PriceMax, err = strconv.Atoi(priceMaxStr)
	if err != nil {
		options.PriceMax = 0
	}

	bedsMaxStr := r.URL.Query().Get("bedsMax")
	options.BedsMax, err = strconv.Atoi(bedsMaxStr)
	if err != nil {
		options.BedsMax = 0
	}

	bathsMaxStr := r.URL.Query().Get("bathsMax")
	options.BathsMax, err = strconv.Atoi(bathsMaxStr)
	if err != nil {
		options.BathsMax = 0
	}

	squareFootageMinStr := r.URL.Query().Get("squareFootageMin")
	options.SquareFootageMin, err = strconv.Atoi(squareFootageMinStr)
	if err != nil {
		options.SquareFootageMin = 0
	}

	squareFootageMaxStr := r.URL.Query().Get("squareFootageMax")
	options.SquareFootageMax, err = strconv.Atoi(squareFootageMaxStr)
	if err != nil {
		options.SquareFootageMax = 0
	}

	lotSizeMinStr := r.URL.Query().Get("lotSizeMin")
	options.LotSizeMin, err = strconv.Atoi(lotSizeMinStr)
	if err != nil {
		options.LotSizeMin = 0
	}

	lotSizeMaxStr := r.URL.Query().Get("lotSizeMax")
	options.LotSizeMax, err = strconv.Atoi(lotSizeMaxStr)
	if err != nil {
		options.LotSizeMax = 0
	}

	yearBuiltMinStr := r.URL.Query().Get("yearBuiltMin")
	options.YearBuiltMin, err = strconv.Atoi(yearBuiltMinStr)
	if err != nil {
		options.YearBuiltMin = 0
	}

	yearBuiltMaxStr := r.URL.Query().Get("yearBuiltMax")
	options.YearBuiltMax, err = strconv.Atoi(yearBuiltMaxStr)
	if err != nil {
		options.YearBuiltMax = 0
	}

	garageSpacesMinStr := r.URL.Query().Get("garageSpacesMin")
	options.GarageSpacesMin, err = strconv.Atoi(garageSpacesMinStr)
	if err != nil {
		options.GarageSpacesMin = 0
	}

	garageSpacesMaxStr := r.URL.Query().Get("garageSpacesMax")
	options.GarageSpacesMax, err = strconv.Atoi(garageSpacesMaxStr)
	if err != nil {
		options.GarageSpacesMax = 0
	}

	propertyTypesStr := r.URL.Query().Get("propertyTypes")
	if err = json.Unmarshal([]byte(propertyTypesStr), &options.PropertyTypes); err != nil {
		options.PropertyTypes = nil
	}

	coolingStr := r.URL.Query().Get("cooling")
	options.Cooling, err = strconv.Atoi(coolingStr)
	if err != nil {
		options.Cooling = 0
	}

	heatingStr := r.URL.Query().Get("heating")
	options.Heating, err = strconv.Atoi(heatingStr)
	if err != nil {
		options.Heating = 0
	}

	poolStr := r.URL.Query().Get("pool")
	options.Pool, err = strconv.Atoi(poolStr)
	if err != nil {
		options.Pool = 0
	}

	return options
}
