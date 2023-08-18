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

func NewRepo(params Params) interfaces.CommentRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) GetLendingByID(ctx context.Context, id int) (data structs.Lending, err error) {
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
    images,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
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
		&data.FeaturesAndAmenities,
		&data.Title,
		&data.Description,
		&data.Video,
		&data.Images,
		&data.CreatedAt,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.GetLendingByID r.db.QueryRow", zap.Error(err))
		return data, err
	}

	return data, nil
}

func (r *repo) GetLendingList(ctx context.Context) (list []structs.LendingList, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    name
FROM lending
WHERE 1=1;`)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.GetLendingList r.db.Query", zap.Error(err))
		return list, err
	}

	for rows.Next() {
		var data structs.LendingList
		err = rows.Scan(
			&data.ID,
			&data.Name,
		)
		if err != nil {
			r.logger.Error("pkg.repo.admin.lending.GetLendingList rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, data)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	return list, nil
}

func (r *repo) SaveLending(ctx context.Context, leadID int) (comments []structs.Comment, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    lead_id,
    admin_id, 
    admin_login, 
    text, 
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM comment
WHERE lead_id = $1
ORDER BY id DESC;`, leadID)
	if err != nil {
		r.logger.Error("pkg.repo.comment.GetCommentsByLeadID r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var c structs.Comment
		err = rows.Scan(
			&c.ID,
			&c.LeadID,
			&c.AdminID,
			&c.AdminLogin,
			&c.Text,
			&c.CreatedAt,
		)
		if err != nil {
			r.logger.Error("pkg.repo.comment.GetCommentsByLeadID rows.Scan()", zap.Error(err))
			return nil, err
		}

		comments = append(comments, c)
	}

	if len(comments) == 0 {
		return nil, errors.ErrNotFound
	}

	r.logger.Info("pkg.repo.comment.GetCommentsByLeadID ", zap.Any("comments", comments))

	return comments, nil
}
