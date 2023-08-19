package lending

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

func NewRepo(params Params) interfaces.LendingClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SelectLendingData(ctx context.Context, id int) (data structs.Lending, err error) {
	err = r.db.QueryRow(ctx, `
SELECT 
    id,
    name,
    full_name,
    address,
    starting_price_aed,
    starting_price_usd,
    property_type,
    furnishing,
    features_and_amenities,
    title,
    description,
    video,
    images
FROM lending
WHERE id = $1;`, id).Scan(
		&data.ID,
		&data.Name,
		&data.FullName,
		&data.Address,
		&data.StartingPrice.AED,
		&data.StartingPrice.USD,
		&data.ListingDetails.PropertyType,
		&data.ListingDetails.Furnishing,
		&data.FeaturesAndAmenitiesIDs,
		&data.Title,
		&data.Description,
		&data.Video,
		&data.Images,
	)
	if err != nil {
		r.logger.Error("pkg.repo.client.lending.SelectLendingData r.db.QueryRow", zap.Error(err))
		return data, err
	}

	return data, nil
}
