package handlers

import (
	"crm/internal/structs"
	"encoding/json"
	"net/http"
	"strconv"
)

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

	bedsMinStr := r.URL.Query().Get("bedsMin")
	options.BedsMin, err = strconv.Atoi(bedsMinStr)
	if err != nil {
		options.BedsMin = 0
	}

	bathsMinStr := r.URL.Query().Get("bathsMin")
	options.BathsMin, err = strconv.Atoi(bathsMinStr)
	if err != nil {
		options.BathsMin = 0
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
