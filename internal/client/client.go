package client

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewService)

type Params struct {
	fx.In
	Logger         *zap.Logger
	LeadRepo       interfaces.LeadClientRepo
	EstateRepo     interfaces.EstateClientRepo
	TextRepo       interfaces.TextClientRepo
	FavoritesRepo  interfaces.FavoritesClientRepo
	LendingRepo    interfaces.LendingClientRepo
	ServiceKeyRepo interfaces.ServiceKeysClientRepo
}

type Service interface {
	SaveLead(ctx context.Context, request structs.Lead) error

	GetEstates(ctx context.Context, offset, limit int, options structs.SearchOptionsDTO) (estates []structs.EstateForList, totalCount int, err error)
	GetLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, totalCount int, err error)
	GetEstateByID(ctx context.Context, id int) (estate structs.Estate, err error)
	GetSearchOptions(ctx context.Context) (option structs.SearchOptions, err error)
	ClearCache()
	GetImageBaseURL() string

	GetTexts(ctx context.Context) (texts []structs.Text, err error)

	SaveFavorite(ctx context.Context, favorite structs.Favorite) error
	DeleteFavorite(ctx context.Context, favorite structs.Favorite) error
	GetEstateFavorites(ctx context.Context, userID int) (estates []structs.EstateForList, err error)

	GetLendingData(ctx context.Context, key string) (data structs.LendingData, err error)
	GetLandingsList(ctx context.Context) (data structs.ListMainPage, err error)
}

type service struct {
	logger         *zap.Logger
	leadRepo       interfaces.LeadClientRepo
	estateRepo     interfaces.EstateClientRepo
	textRepo       interfaces.TextClientRepo
	favoritesRepo  interfaces.FavoritesClientRepo
	lendingRepo    interfaces.LendingClientRepo
	serviceKeyRepo interfaces.ServiceKeysClientRepo
}

func NewService(params Params) Service {
	return &service{
		logger:         params.Logger,
		leadRepo:       params.LeadRepo,
		estateRepo:     params.EstateRepo,
		textRepo:       params.TextRepo,
		favoritesRepo:  params.FavoritesRepo,
		lendingRepo:    params.LendingRepo,
		serviceKeyRepo: params.ServiceKeyRepo,
	}
}
