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
	Logger     *zap.Logger
	LeadRepo   interfaces.LeadClientRepo
	EstateRepo interfaces.EstateClientRepo
}

type Service interface {
	SaveLead(ctx context.Context, request structs.Lead) error

	GetEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error)
	GetLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error)
	GetEstateByID(ctx context.Context, id int) (estate structs.Estate, err error)
	GetSearchOptions(ctx context.Context) (option structs.SearchOptions, err error)
	ClearCache()
}

type service struct {
	logger     *zap.Logger
	leadRepo   interfaces.LeadClientRepo
	estateRepo interfaces.EstateClientRepo
}

func NewService(params Params) Service {
	return &service{
		logger:     params.Logger,
		leadRepo:   params.LeadRepo,
		estateRepo: params.EstateRepo,
	}
}
