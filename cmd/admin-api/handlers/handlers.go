package handlers

import (
	"crm/internal/admin"
	"crm/internal/auth"
	"crm/internal/middleware"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
)

var Module = fx.Provide(NewHandler)

type Handler interface {
	SignIn(http.ResponseWriter, *http.Request)
	GetLeadList(http.ResponseWriter, *http.Request)
	GetLeadListXlSX(http.ResponseWriter, *http.Request)
	GetLeadAndCommentsByID(http.ResponseWriter, *http.Request)
	AddComment(http.ResponseWriter, *http.Request)
	UpdateLeadStatus(http.ResponseWriter, *http.Request)

	GetEstateByID(http.ResponseWriter, *http.Request)
	GetEstates(http.ResponseWriter, *http.Request)
	AddEstate(http.ResponseWriter, *http.Request)
	UpdateEstate(http.ResponseWriter, *http.Request)
	DeleteEstate(http.ResponseWriter, *http.Request)
	ApprovedEstate(http.ResponseWriter, *http.Request)
	UploadImages(http.ResponseWriter, *http.Request)
	DeleteImage(http.ResponseWriter, *http.Request)
	GetImageBaseURL(http.ResponseWriter, *http.Request)

	GetTexts(http.ResponseWriter, *http.Request)
	UpdateText(http.ResponseWriter, *http.Request)

	MwCheckAdminAuthToken(m middleware.Handler) middleware.Handler
}

type HandlerParams struct {
	fx.In
	Logger       *zap.Logger
	AdminService admin.Service
	AuthService  auth.Service
}

type handler struct {
	logger       *zap.Logger
	adminService admin.Service
	authService  auth.Service
}

func NewHandler(params HandlerParams) Handler {
	return &handler{
		logger:       params.Logger,
		adminService: params.AdminService,
		authService:  params.AuthService,
	}
}
