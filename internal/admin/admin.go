package admin

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"mime/multipart"
)

var Module = fx.Provide(NewService)

type Params struct {
	fx.In
	Logger      *zap.Logger
	LeadRepo    interfaces.LeadAdminRepo
	CommentRepo interfaces.CommentRepo
	EstateRepo  interfaces.EstateAdminRepo
	TextRepo    interfaces.TextAdminRepo
	LandingRepo interfaces.LandingAdminRepo
}

type Service interface {
	// lead
	GetLeadList(ctx context.Context, offset, limit int, status string) (list []structs.Lead, err error)
	GetLeadListXLSX(ctx context.Context, status string) (url string, err error)
	GetLeadAndCommentsByID(ctx context.Context, id int) (lead structs.LeadAndComments, err error)
	AddComment(ctx context.Context, comment structs.Comment) error
	UpdateLeadStatus(ctx context.Context, id int, status string) error

	// home website
	GetEstateByID(ctx context.Context, id int) (structs.Estate, error)
	GetEstates(ctx context.Context, offset, limit int, status string) ([]structs.EstateForList, int, error)
	AddEstate(ctx context.Context, request structs.Estate) error
	UpdateEstate(ctx context.Context, request structs.Estate) error
	DeleteEstate(ctx context.Context, id int) error
	ApprovedEstate(ctx context.Context, id int) error
	UploadEstateImages(ctx context.Context, id int, files []multipart.File) error
	DeleteEstateImages(ctx context.Context, id int, imageName string) error
	GetImageBaseURL() string

	GetClientsFavorites(ctx context.Context, offset, limit int) (clientsFavorites []structs.Client, err error)

	GetTexts(ctx context.Context) (texts []structs.Text, err error)
	UpdateText(ctx context.Context, text structs.Text) error

	//landing
	SaveLanding(ctx context.Context, data structs.Landing) error
	UpdateLanding(ctx context.Context, data structs.Landing) error
	GetLandingData(ctx context.Context, landingID int) (data structs.Landing, err error)
	GetLandingList(ctx context.Context, offset, limit int) (list []structs.LandingList, count int, err error)
	UploadLandingImages(ctx context.Context, id int, files []multipart.File) error
	DeleteLandingImages(ctx context.Context, id int, imageName string) error
	UploadBackgroundImage(ctx context.Context, landingID int, file multipart.File, typeName string) error
	UploadBackgroundForMobile(ctx context.Context, landingID int, file multipart.File, typeName string) error
	UploadMainLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error
	UploadPartnerLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error
	UploadOurLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error
	UploadFilePlan(ctx context.Context, availabilityID int, file multipart.File, typeName string) error
	UploadVideoCover(ctx context.Context, landingID int, file multipart.File, typeName string) error
	GetFileURL(ctx context.Context) string
	GetSpecialGiftIcons(ctx context.Context) []structs.SpecialGiftIcon

	AddFeatureAndAmenity(ctx context.Context, file multipart.File, typeName, featureName string) error
	DeleteFeatureAndAmenity(ctx context.Context, id int) error
	GetFeaturesAndAmenities(ctx context.Context) (list []structs.FeatureOrAmenity, err error)

	SaveAvailability(ctx context.Context, data structs.Availability) error
	UpdateAvailability(ctx context.Context, data structs.Availability) error
	RemoveAvailability(ctx context.Context, id int) error
}

type service struct {
	logger      *zap.Logger
	leadRepo    interfaces.LeadAdminRepo
	commentRepo interfaces.CommentRepo
	estateRepo  interfaces.EstateAdminRepo
	textRepo    interfaces.TextAdminRepo
	landingRepo interfaces.LandingAdminRepo
}

func NewService(params Params) Service {
	return &service{
		logger:      params.Logger,
		leadRepo:    params.LeadRepo,
		commentRepo: params.CommentRepo,
		estateRepo:  params.EstateRepo,
		textRepo:    params.TextRepo,
		landingRepo: params.LandingRepo,
	}
}
