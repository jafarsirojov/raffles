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

func NewRepo(params Params) interfaces.LandingClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SelectLandingList(ctx context.Context) (landings []structs.LandingListMainPage, err error) {
	rows, err := r.db.Query(ctx, `
SELECT l.name,
       l.main_description,
       l.background_image,
       sk.key
FROM landing l
         JOIN service_keys sk on l.id = sk.landing_id
WHERE 1 = 1;`)
	if err != nil {
		r.logger.Error("pkg.repo.client.landing.SelectLandingList r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var l structs.LandingListMainPage
		err = rows.Scan(
			&l.Name,
			&l.MainDescription,
			&l.BackgroundImage,
			&l.Key,
		)
		if err != nil {
			r.logger.Error("pkg.repo.client.landing.SelectLandingList rows.Scan", zap.Error(err))
			return nil, err
		}

		landings = append(landings, l)
	}

	if len(landings) == 0 {
		return nil, errors.ErrNotFound
	}

	return landings, nil
}

func (r *repo) SelectLandingData(ctx context.Context, id int) (data structs.Landing, err error) {
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
	location_description
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
	)
	if err != nil {
		r.logger.Error("pkg.repo.client.landing.SelectLandingData r.db.QueryRow", zap.Error(err))
		return data, err
	}

	return data, nil
}
