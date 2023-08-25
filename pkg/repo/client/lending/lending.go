package lending

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

func (r *repo) SelectLandingList(ctx context.Context) (landings []structs.LendingListMainPage, err error) {
	rows, err := r.db.Query(ctx, `
SELECT l.name,
       l.main_description,
       l.background_image,
       sk.key
FROM lending l
         JOIN service_keys sk on l.id = sk.lending_id
WHERE 1 = 1;`)
	if err != nil {
		r.logger.Error("pkg.repo.client.lending.SelectLandingList r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var l structs.LendingListMainPage
		err = rows.Scan(
			&l.Name,
			&l.MainDescription,
			&l.BackgroundImage,
			&l.Key,
		)
		if err != nil {
			r.logger.Error("pkg.repo.client.lending.SelectLandingList rows.Scan", zap.Error(err))
			return nil, err
		}

		landings = append(landings, l)
	}

	if len(landings) == 0 {
		return nil, errors.ErrNotFound
	}

	return landings, nil
}

func (r *repo) SelectLendingData(ctx context.Context, id int) (data structs.Lending, err error) {
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
    file_plan,
    title_plan,
    images,
    background_image,
    main_logo,
	partner_logo, 
	our_logo
FROM lending
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
		&data.FilePlan,
		&data.TitlePlan,
		&data.Images,
		&data.BackgroundImage,
		&data.MainLogo,
		&data.PartnerLogo,
		&data.OurLogo,
	)
	if err != nil {
		r.logger.Error("pkg.repo.client.lending.SelectLendingData r.db.QueryRow", zap.Error(err))
		return data, err
	}

	return data, nil
}
