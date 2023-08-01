package handlers

import (
	"crm/internal/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Provide(NewHandler)

type Handler interface {
	SaveLead(http.ResponseWriter, *http.Request)

	GetEstates(http.ResponseWriter, *http.Request)
	GetLuxuryEstates(http.ResponseWriter, *http.Request)
	GetEstateByID(http.ResponseWriter, *http.Request)

	GetSearchOptions(http.ResponseWriter, *http.Request)
	ClearCache(http.ResponseWriter, *http.Request)
	GetImageBaseURL(http.ResponseWriter, *http.Request)

	GetTexts(http.ResponseWriter, *http.Request)
}

type HandlerParams struct {
	fx.In
	Logger        *zap.Logger
	ClientService client.Service
}

type handler struct {
	logger        *zap.Logger
	clientService client.Service
}

func NewHandler(params HandlerParams) Handler {
	return &handler{
		logger:        params.Logger,
		clientService: params.ClientService,
	}
}
