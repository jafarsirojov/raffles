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
	LendingRepo interfaces.LendingAdminRepo
}

type Service interface {
	GetLeadList(ctx context.Context, offset, limit int, status string) (list []structs.Lead, err error)
	GetLeadListXLSX(ctx context.Context, status string) (url string, err error)
	GetLeadAndCommentsByID(ctx context.Context, id int) (lead structs.LeadAndComments, err error)
	AddComment(ctx context.Context, comment structs.Comment) error
	UpdateLeadStatus(ctx context.Context, id int, status string) error

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

	//lending
	SaveLending(ctx context.Context, data structs.Lending) error
	UpdateLending(ctx context.Context, data structs.Lending) error
	GetLendingData(ctx context.Context, landingID int) (data structs.Lending, err error)
	GetLendingList(ctx context.Context) (list []structs.LendingList, err error)
	UploadLendingImages(ctx context.Context, id int, files []multipart.File) error
	DeleteLendingImages(ctx context.Context, id int, imageName string) error
}

type service struct {
	logger      *zap.Logger
	leadRepo    interfaces.LeadAdminRepo
	commentRepo interfaces.CommentRepo
	estateRepo  interfaces.EstateAdminRepo
	textRepo    interfaces.TextAdminRepo
	lendingRepo interfaces.LendingAdminRepo
}

func NewService(params Params) Service {
	return &service{
		logger:      params.Logger,
		leadRepo:    params.LeadRepo,
		commentRepo: params.CommentRepo,
		estateRepo:  params.EstateRepo,
		textRepo:    params.TextRepo,
	}
}
