package landing

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/db"
	"crm/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewRepo)

type Params struct {
	fx.In
	DB     db.Querier
	Logger *zap.Logger
}

func NewRepo(params Params) interfaces.LandingAdminRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SaveLanding(ctx context.Context, data structs.Landing) (err error) {
	_, err = r.db.Exec(ctx, `
INSERT INTO landing(
    name,
    main_description,
    full_name,
    slogan,
    address,
    starting_price_aed,
    starting_price_usd,
    property_type,
    furnishing,
    features_and_amenities,
    title,
    description,
    video,
    title_plan,
    latitude,
    longitude,
    location_description,
    images) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18);`,
		data.Name,
		data.MainDescription,
		data.FullName,
		data.Slogan,
		data.Address,
		data.StartingPrice.AED,
		data.StartingPrice.USD,
		data.ListingDetails.PropertyType,
		data.ListingDetails.Furnishing,
		data.FeaturesAndAmenitiesIDs,
		data.Title,
		data.Description,
		data.Video,
		data.TitlePlan,
		data.Latitude,
		data.Longitude,
		data.LocationDescription,
		[]string{},
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SaveLanding r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) GetLandingByID(ctx context.Context, id int) (data structs.Landing, err error) {
	err = r.db.QueryRow(ctx, `
SELECT 
    id,
    name,
    main_description,
    full_name,
    slogan,
    address,
    starting_price_aed,
    starting_price_usd,
    property_type,
    furnishing,
    features_and_amenities,
    title,
    description,
    video,
    video_cover,
    file_plan,
    title_plan, 
    images,
    background_image,
    background_for_mobile,
	main_logo,
	partner_logo, 
	our_logo,
	latitude,
    longitude,
    location_description,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI'),
    to_char(updated_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM landing
WHERE id = $1;`, id).Scan(
		&data.ID,
		&data.Name,
		&data.MainDescription,
		&data.FullName,
		&data.Slogan,
		&data.Address,
		&data.StartingPrice.AED,
		&data.StartingPrice.USD,
		&data.ListingDetails.PropertyType,
		&data.ListingDetails.Furnishing,
		&data.FeaturesAndAmenitiesIDs,
		&data.Title,
		&data.Description,
		&data.Video,
		&data.VideoCover,
		&data.FilePlan,
		&data.TitlePlan,
		&data.Images,
		&data.BackgroundImage,
		&data.BackgroundForMobile,
		&data.MainLogo,
		&data.PartnerLogo,
		&data.OurLogo,
		&data.Latitude,
		&data.Longitude,
		&data.LocationDescription,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.GetLandingByID r.db.QueryRow", zap.Error(err))
		return data, err
	}

	return data, nil
}

func (r *repo) GetLandingList(ctx context.Context, offset, limit int) (list []structs.LandingList, count int, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    name
FROM landing
WHERE 1=1
ORDER BY id DESC OFFSET $1 LIMIT $2;`, offset, limit)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.GetLandingList r.db.Query", zap.Error(err))
		return list, 0, err
	}

	for rows.Next() {
		var data structs.LandingList
		err = rows.Scan(
			&data.ID,
			&data.Name,
		)
		if err != nil {
			r.logger.Error("pkg.repo.admin.landing.GetLandingList rows.Scan()", zap.Error(err))
			return nil, 0, err
		}

		list = append(list, data)
	}

	if len(list) == 0 {
		return nil, 0, errors.ErrNotFound
	}

	err = r.db.QueryRow(ctx, "SELECT count(1) FROM landing WHERE 1=1;").Scan(&count)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.GetLandingList rows.Scan()", zap.Error(err))
	}

	return list, count, nil
}

func (r *repo) UpdateLanding(ctx context.Context, data structs.Landing) (err error) {
	_, err = r.db.Exec(ctx, `
UPDATE landing SET
    name = $2,
    full_name = $3,
    address = $4,
    starting_price_aed = $5,
    starting_price_usd = $6,
    property_type = $7,
    furnishing = $8,
    features_and_amenities = $9,
    title = $10,
    description = $11,
    video = $12,
    main_description = $13,
    slogan = $14,
    title_plan = $15,
    latitude = $16,
    longitude = $17,
    location_description = $18,
    updated_at = now()
    WHERE id = $1`,
		data.ID,
		data.Name,
		data.FullName,
		data.Address,
		data.StartingPrice.AED,
		data.StartingPrice.USD,
		data.ListingDetails.PropertyType,
		data.ListingDetails.Furnishing,
		data.FeaturesAndAmenitiesIDs,
		data.Title,
		data.Description,
		data.Video,
		data.MainDescription,
		data.Slogan,
		data.TitlePlan,
		data.Latitude,
		data.Longitude,
		data.LocationDescription,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateLanding r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) SelectFilePlanByLandingID(ctx context.Context, id int) (paymentPlan string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT file_plan FROM landing WHERE id = $1;`, id).Scan(&paymentPlan)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectFilePlanByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return paymentPlan, err
	}

	return paymentPlan, nil
}

func (r *repo) UpdateFilePlan(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET file_plan = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateFilePlan r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}
