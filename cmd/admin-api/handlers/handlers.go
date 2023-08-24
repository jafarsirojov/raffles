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
	UploadEstateImages(http.ResponseWriter, *http.Request)
	DeleteEstateImages(http.ResponseWriter, *http.Request)
	GetImageBaseURL(http.ResponseWriter, *http.Request)

	GetClientsFavorites(http.ResponseWriter, *http.Request)

	GetTexts(http.ResponseWriter, *http.Request)
	UpdateText(http.ResponseWriter, *http.Request)

	// lending
	AddLendingPage(http.ResponseWriter, *http.Request)
	UpdateLendingPage(http.ResponseWriter, *http.Request)
	GetLendingData(http.ResponseWriter, *http.Request)
	GetLendingList(http.ResponseWriter, *http.Request)
	UploadLendingImages(http.ResponseWriter, *http.Request)
	DeleteLendingImages(http.ResponseWriter, *http.Request)
	UploadBackgroundImage(http.ResponseWriter, *http.Request)
	UploadMainLogo(http.ResponseWriter, *http.Request)
	UploadPartnerLogo(http.ResponseWriter, *http.Request)
	UploadOurLogo(http.ResponseWriter, *http.Request)
	UploadFilePlan(http.ResponseWriter, *http.Request)

	GetFeaturesAndAmenities(http.ResponseWriter, *http.Request)
	AddFeatureAndAmenity(http.ResponseWriter, *http.Request)
	DeleteFeatureAndAmenity(http.ResponseWriter, *http.Request)

	UpdateAvailability(http.ResponseWriter, *http.Request)
	AddAvailability(http.ResponseWriter, *http.Request)
	RemoveAvailability(http.ResponseWriter, *http.Request)

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
