package handlers

import (
	"crm/internal/authClient"
	"crm/internal/client"
	"crm/internal/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Provide(NewHandler)

type Handler interface {
	SaveLead(http.ResponseWriter, *http.Request)

	// Home Site

	GetEstates(http.ResponseWriter, *http.Request)
	GetLuxuryEstates(http.ResponseWriter, *http.Request)
	GetEstateByID(http.ResponseWriter, *http.Request)

	GetSearchOptions(http.ResponseWriter, *http.Request)
	ClearCache(http.ResponseWriter, *http.Request)
	GetImageBaseURL(http.ResponseWriter, *http.Request)

	GetTexts(http.ResponseWriter, *http.Request)

	SignUp(http.ResponseWriter, *http.Request)
	SignIn(http.ResponseWriter, *http.Request)

	SaveFavorite(http.ResponseWriter, *http.Request)
	DeleteFavorite(http.ResponseWriter, *http.Request)
	GetFavorites(http.ResponseWriter, *http.Request)
	GetLendingData(http.ResponseWriter, *http.Request)

	MwCheckAuthToken(middleware middleware.Handler) middleware.Handler
}

type HandlerParams struct {
	fx.In
	Logger        *zap.Logger
	ClientService client.Service
	AuthService   authClient.Service
}

type handler struct {
	logger        *zap.Logger
	clientService client.Service
	authService   authClient.Service
}

func NewHandler(params HandlerParams) Handler {
	return &handler{
		logger:        params.Logger,
		clientService: params.ClientService,
		authService:   params.AuthService,
	}
}
