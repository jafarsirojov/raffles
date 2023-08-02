package estate

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/db"
	"crm/pkg/errors"
	"github.com/jackc/pgx/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB     db.Querier
	Logger *zap.Logger
}

func NewRepo(params Params) interfaces.EstateClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SelectEstateByID(ctx context.Context, id int) (estate structs.Estate, err error) {
	err = r.db.QueryRow(ctx, `
SELECT 
    id,
    status,
    luxury,
    images,
    
    name,
    price, 
    country, 
    city, 
    address, 
    beds,
    baths,
    area_in_meter,
    property_type,
    year_built,
    description,
    latitude,
    longitude,
    
    appliances,
    interior_features,
    kitchen_features,
    total_bedrooms,
    full_bathrooms,
	half_bathrooms,
	floor_description,
	fireplace,
	cooling,
	heating,
    
    lot_size_in_acres,
	exterior_features,
	arch_style,
	roof,
	sewer,
    
    area_name,
	garage,
	parking,
	view,
	pool,
	pool_description,
	water_source,
	utilities,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI'), 
  	to_char(updated_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM estate
WHERE id = $1 and status = 'active';`, id).Scan(
		&estate.ID,
		&estate.Status,
		&estate.Luxury,
		&estate.Images,

		&estate.PropertyDescription.Name,
		&estate.PropertyDescription.Price,
		&estate.PropertyDescription.Country,
		&estate.PropertyDescription.City,
		&estate.PropertyDescription.Address,
		&estate.PropertyDescription.Beds,
		&estate.PropertyDescription.Baths,
		&estate.PropertyDescription.AreaInMeter,
		&estate.PropertyDescription.Type,
		&estate.PropertyDescription.YearBuilt,
		&estate.PropertyDescription.Description,
		&estate.PropertyDescription.Latitude,
		&estate.PropertyDescription.Longitude,

		&estate.Interior.Appliances,
		&estate.Interior.Features,
		&estate.Interior.KitchenFeatures,
		&estate.Interior.TotalBedrooms,
		&estate.Interior.FullBathrooms,
		&estate.Interior.HalfBathrooms,
		&estate.Interior.FloorDescription,
		&estate.Interior.Fireplace,
		&estate.Interior.Cooling,
		&estate.Interior.Heating,

		&estate.Exterior.LotSizeInAcres,
		&estate.Exterior.Features,
		&estate.Exterior.ArchStyle,
		&estate.Exterior.Roof,
		&estate.Exterior.Sewer,

		&estate.OtherDetails.AreaName,
		&estate.OtherDetails.Garage,
		&estate.OtherDetails.Parking,
		&estate.OtherDetails.View,
		&estate.OtherDetails.Pool,
		&estate.OtherDetails.PoolDescription,
		&estate.OtherDetails.WaterSource,
		&estate.OtherDetails.Utilities,

		&estate.CreateAt,
		&estate.UpdateAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return estate, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.client.estate.GetEstateByID r.db.Query", zap.Error(err))
		return estate, err
	}

	r.logger.Info("pkg.repo.client.estate.GetEstateByID", zap.Any("estate", estate))

	return estate, nil
}

func (r *repo) SelectEstates(ctx context.Context, offset, limit int, options structs.SearchOptionsDTO) (estates []structs.EstateForList, err error) {
	where, fields := r.getWhereByOptions(options, offset, limit)

	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    status,
    images,
    
    name,
    price, 
    address, 
    beds,
    baths,
    area_in_meter,
    latitude,
    longitude
FROM estate
WHERE status = 'active'`+where, fields...)
	if err != nil {
		r.logger.Error("pkg.repo.client.estate.GetEstates r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var estate structs.EstateForList
		err = rows.Scan(
			&estate.ID,
			&estate.Status,
			&estate.Images,

			&estate.Name,
			&estate.Price,
			&estate.Address,
			&estate.Beds,
			&estate.Baths,
			&estate.AreaInMeter,
			&estate.Latitude,
			&estate.Longitude,
		)
		if err != nil {
			r.logger.Error("pkg.repo.estate.GetEstates rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) SelectLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    status,
    images,
    
    name,
    price, 
    address, 
    beds,
    baths,
    area_in_meter,
    latitude,
    longitude
FROM estate
WHERE status = 'active' and luxury = true
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit)
	if err != nil {
		r.logger.Error("pkg.repo.client.estate.SelectLuxuryEstates r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var estate structs.EstateForList
		err = rows.Scan(
			&estate.ID,
			&estate.Status,
			&estate.Images,

			&estate.Name,
			&estate.Price,
			&estate.Address,
			&estate.Beds,
			&estate.Baths,
			&estate.AreaInMeter,
			&estate.Latitude,
			&estate.Longitude,
		)
		if err != nil {
			r.logger.Error("pkg.repo.estate.SelectLuxuryEstates rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) SelectSearchOptions(ctx context.Context) (options structs.SearchOptions, err error) {
	err = r.db.QueryRow(ctx, `
SELECT min(price),
       max(price),
       max(beds),
       max(baths),
       min(area_in_meter),
       max(area_in_meter),
       min(lot_size_in_acres),
       max(lot_size_in_acres),
       min(year_built),
       max(year_built),
       min(garage),
       max(garage)
from estate
WHERE status = 'active';`,
	).Scan(
		&options.PriceMin,
		&options.PriceMax,
		&options.BedsMax,
		&options.BathsMax,
		&options.SquareFootageMin,
		&options.SquareFootageMax,
		&options.LotSizeMin,
		&options.LotSizeMax,
		&options.YearBuiltMin,
		&options.YearBuiltMax,
		&options.GarageSpacesMin,
		&options.GarageSpacesMax,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return options, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.client.estate.SelectSearchOptions r.db.QueryRow", zap.Error(err))
		return options, err
	}

	return options, nil
}

func (r *repo) GetEstatesTotalCount(ctx context.Context) (count int, err error) {
	err = r.db.QueryRow(ctx, `SELECT count(1) FROM estate WHERE status = 'active';`).Scan(&count)
	if err != nil {
		r.logger.Error("pkg.repo.client.estate.GetEstatesTotalCount r.db.Query", zap.Error(err))
		return 0, err
	}

	return count, nil
}
