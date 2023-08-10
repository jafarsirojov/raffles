package estate

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/db"
	"crm/pkg/errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"strings"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB     db.Querier
	Logger *zap.Logger
}

func NewRepo(params Params) interfaces.EstateAdminRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SaveEstate(ctx context.Context, request structs.Estate) error {
	_, err := r.db.Exec(ctx, `
insert into estate (
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
				  utilities
          		)values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
          		        $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, 
          		        $21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
          		        $31, $32, $33, $34, $35, $36, $37, $38, $39);`,
		request.Status,
		request.Luxury,
		[]string{},

		request.PropertyDescription.Name,
		request.PropertyDescription.Price,
		request.PropertyDescription.Country,
		request.PropertyDescription.City,
		request.PropertyDescription.Address,
		request.PropertyDescription.Beds,
		request.PropertyDescription.Baths,
		request.PropertyDescription.AreaInMeter,
		request.PropertyDescription.Type,
		request.PropertyDescription.YearBuilt,
		request.PropertyDescription.Description,
		request.PropertyDescription.Latitude,
		request.PropertyDescription.Longitude,

		request.Interior.Appliances,
		request.Interior.Features,
		request.Interior.KitchenFeatures,
		request.Interior.TotalBedrooms,
		request.Interior.FullBathrooms,
		request.Interior.HalfBathrooms,
		request.Interior.FloorDescription,
		request.Interior.Fireplace,
		request.Interior.Cooling,
		request.Interior.Heating,

		request.Exterior.LotSizeInAcres,
		request.Exterior.Features,
		request.Exterior.ArchStyle,
		request.Exterior.Roof,
		request.Exterior.Sewer,

		request.OtherDetails.AreaName,
		request.OtherDetails.Garage,
		request.OtherDetails.Parking,
		request.OtherDetails.View,
		request.OtherDetails.Pool,
		request.OtherDetails.PoolDescription,
		request.OtherDetails.WaterSource,
		request.OtherDetails.Utilities,
	)
	if err != nil {
		r.logger.Error("pkg.admin.repo.estate.AddEstate r.db.Exec", zap.Error(err), zap.Any("request", request))
		return err
	}

	return nil
}

func (r *repo) GetEstateByID(ctx context.Context, id int) (estate structs.Estate, err error) {
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
WHERE id = $1 and status <> 'deleted';`, id).Scan(
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
		r.logger.Error("pkg.admin.repo.estate.GetEstateByID r.db.Query", zap.Error(err))
		return estate, err
	}

	return estate, nil
}

func (r *repo) GetEstates(ctx context.Context, offset, limit int) (estates []structs.Estate, err error) {
	rows, err := r.db.Query(ctx, `
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
WHERE status <> 'deleted'
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit)
	if err != nil {
		r.logger.Error("pkg.repo.admin.estate.GetEstates r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var estate structs.Estate
		err = rows.Scan(
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
			r.logger.Error("pkg.repo.admin.estate.GetEstates rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) GetEstatesForList(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error) {
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
WHERE status <> 'deleted'
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit)
	if err != nil {
		r.logger.Error("pkg.repo.admin.estate.GetEstatesForList r.db.Query", zap.Error(err))
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
			r.logger.Error("pkg.repo.admin.estate.GetEstatesForList rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) GetEstatesForListByStatus(ctx context.Context, offset, limit int, status string) (estates []structs.EstateForList, err error) {
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
WHERE status = $3
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit, status)
	if err != nil {
		r.logger.Error("pkg.repo.admin.estate.GetEstatesForListByStatus r.db.Query", zap.Error(err))
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
			r.logger.Error("pkg.repo.admin.estate.GetEstatesForListByStatus rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) GetEstatesByStatus(ctx context.Context, offset, limit int, status string) (estates []structs.Estate, err error) {
	rows, err := r.db.Query(ctx, `
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
WHERE status = $3
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit, status)
	if err != nil {
		r.logger.Error("pkg.repo.admin.estate.GetEstatesByStatus r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var estate structs.Estate
		err = rows.Scan(
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
			r.logger.Error("pkg.repo.admin.estate.GetEstatesByStatus rows.Scan()", zap.Error(err))
			return nil, err
		}

		estates = append(estates, estate)
	}

	if len(estates) == 0 {
		return nil, errors.ErrNotFound
	}

	return estates, nil
}

func (r *repo) GetEstatesTotalCount(ctx context.Context, status string) (count int, err error) {
	var param []interface{}
	c := ""
	if len(strings.TrimSpace(status)) != 0 {
		param = append(param, status)
		c = "AND status = $1"
	}
	err = r.db.QueryRow(ctx, fmt.Sprintf(`SELECT count(1) FROM estate WHERE status <> 'deleted' %s;`, c), param...).Scan(&count)
	if err != nil {
		r.logger.Error("pkg.repo.admin.estate.GetEstatesTotalCount r.db.Query", zap.Error(err))
		return 0, err
	}

	return count, nil
}

func (r *repo) UpdateEstate(ctx context.Context, request structs.Estate) error {
	_, err := r.db.Exec(ctx, `UPDATE estate SET 
    status = $2,
    luxury = $3,
    
    name = $4,
    price = $5,
    country = $6,
    city = $7,
    address = $8,
    beds = $9,
    baths = $10,
    area_in_meter = $11,
    property_type = $12,
    year_built = $13,
    description = $14,
    latitude = $15,
    longitude = $16,
    
    appliances = $17,
    interior_features = $18,
    kitchen_features = $19,
    total_bedrooms = $20,
    full_bathrooms = $21,
	half_bathrooms = $22,
	floor_description = $23,
	fireplace = $24,
	cooling = $25,
	heating = $26,
    
    lot_size_in_acres = $27,
	exterior_features = $28,
	arch_style = $29,
	roof = $30,
	sewer = $31,
    
    area_name = $32,
	garage = $33,
	parking = $34,
	view = $35,
	pool = $36,
	pool_description = $37,
	water_source = $38,
	utilities = $39,
	updated_at = now()
WHERE id = $1;`,
		request.ID,

		request.Status,
		request.Luxury,

		request.PropertyDescription.Name,
		request.PropertyDescription.Price,
		request.PropertyDescription.Country,
		request.PropertyDescription.City,
		request.PropertyDescription.Address,
		request.PropertyDescription.Beds,
		request.PropertyDescription.Baths,
		request.PropertyDescription.AreaInMeter,
		request.PropertyDescription.Type,
		request.PropertyDescription.YearBuilt,
		request.PropertyDescription.Description,
		request.PropertyDescription.Latitude,
		request.PropertyDescription.Longitude,

		request.Interior.Appliances,
		request.Interior.Features,
		request.Interior.KitchenFeatures,
		request.Interior.TotalBedrooms,
		request.Interior.FullBathrooms,
		request.Interior.HalfBathrooms,
		request.Interior.FloorDescription,
		request.Interior.Fireplace,
		request.Interior.Cooling,
		request.Interior.Heating,

		request.Exterior.LotSizeInAcres,
		request.Exterior.Features,
		request.Exterior.ArchStyle,
		request.Exterior.Roof,
		request.Exterior.Sewer,

		request.OtherDetails.AreaName,
		request.OtherDetails.Garage,
		request.OtherDetails.Parking,
		request.OtherDetails.View,
		request.OtherDetails.Pool,
		request.OtherDetails.PoolDescription,
		request.OtherDetails.WaterSource,
		request.OtherDetails.Utilities,
	)
	if err != nil {
		r.logger.Error("pkg.repo.estate.UpdateEstate r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) UpdateEstateStatus(ctx context.Context, id int, status structs.Status) error {
	_, err := r.db.Exec(ctx,
		`UPDATE estate SET status = $2, updated_at = now() WHERE id = $1;`, id, status)
	if err != nil {
		r.logger.Error("pkg.repo.estate.UpdateEstate r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) GetImagesByEstateID(ctx context.Context, id int) (images []string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT images FROM estate WHERE id = $1;`, id).Scan(&images)
	if err != nil {
		r.logger.Error("pkg.repo.estate.GetImagesByEstateID r.db.QueryRow", zap.Error(err))
		return nil, err
	}

	return images, nil
}

func (r *repo) UpdateEstateImages(ctx context.Context, id int, images []string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE estate SET images = $2, updated_at = now() WHERE id = $1;`, id, images)
	if err != nil {
		r.logger.Error("pkg.repo.estate.UpdateEstateImages r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
