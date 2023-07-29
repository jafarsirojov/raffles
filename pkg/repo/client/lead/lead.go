package lead

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/db"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB     db.Querier
	Logger *zap.Logger
}

func NewRepo(params Params) interfaces.LeadClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) InsertLead(ctx context.Context, request structs.Lead) error {
	_, err := r.db.Exec(ctx, `
insert into lead (
          		site,
                name,
          		re_stage_constructor, 
          		re_region, 
          		re_type, 
          		re_purpose_of_acquisition, 
          		re_count_of_rooms,
                purchase_budget,
                phone,
                email,
                communication_method,
                description
          		)values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`,
		request.Site,
		request.Name,
		request.REStageConstruction,
		request.RERegion,
		request.REType,
		request.REPurposeOfAcquisition,
		request.RECountOfRooms,
		request.PurchaseBudget,
		request.Phone,
		request.Email,
		request.CommunicationMethod,
		request.Description,
	)
	if err != nil {
		r.logger.Error("pkg.repo.client.lead.InsertLead r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
