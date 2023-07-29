package lead

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

func NewRepo(params Params) interfaces.LeadAdminRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) GetLeadList(ctx context.Context, offset, limit int) (list []structs.Lead, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
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
    description, 
    status,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI'), 
  	to_char(updated_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM lead
WHERE 1=1 
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit)
	if err != nil {
		r.logger.Error("pkg.repo.lead.GetLeadList r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var l structs.Lead
		err = rows.Scan(
			&l.ID,
			&l.Site,
			&l.Name,
			&l.REStageConstruction,
			&l.RERegion,
			&l.REType,
			&l.REPurposeOfAcquisition,
			&l.RECountOfRooms,
			&l.PurchaseBudget,
			&l.Phone,
			&l.Email,
			&l.CommunicationMethod,
			&l.Description,
			&l.Status,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("pkg.repo.lead.GetLeadList rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, l)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	r.logger.Info("pkg.repo.lead.GetLeadList ", zap.Any("list", list))

	return list, nil
}

func (r *repo) GetLeadListByStatus(ctx context.Context, offset, limit int, status string) (list []structs.Lead, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
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
    description, 
    status,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI'), 
  	to_char(updated_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM lead
WHERE status = $3
ORDER BY id DESC offset $1 limit $2 ;`, offset, limit, status)
	if err != nil {
		r.logger.Error("pkg.repo.lead.GetLeadListByStatus r.db.Query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		var l structs.Lead
		err = rows.Scan(
			&l.ID,
			&l.Site,
			&l.Name,
			&l.REStageConstruction,
			&l.RERegion,
			&l.REType,
			&l.REPurposeOfAcquisition,
			&l.RECountOfRooms,
			&l.PurchaseBudget,
			&l.Phone,
			&l.Email,
			&l.CommunicationMethod,
			&l.Description,
			&l.Status,
			&l.CreatedAt,
			&l.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("pkg.repo.lead.GetLeadListByStatus rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, l)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	r.logger.Info("pkg.repo.lead.GetLeadListByStatus ", zap.Any("list", list))

	return list, nil
}

func (r *repo) GetLeadByID(ctx context.Context, id int) (l structs.Lead, err error) {
	err = r.db.QueryRow(ctx, `
SELECT 
    id,
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
    description, 
    status,
    to_char(created_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI'), 
  	to_char(updated_at AT TIME ZONE 'Asia/Dubai', 'DD-MM HH24:MI')
FROM lead
WHERE id = $1;`, id).Scan(
		&l.ID,
		&l.Site,
		&l.Name,
		&l.REStageConstruction,
		&l.RERegion,
		&l.REType,
		&l.REPurposeOfAcquisition,
		&l.RECountOfRooms,
		&l.PurchaseBudget,
		&l.Phone,
		&l.Email,
		&l.CommunicationMethod,
		&l.Description,
		&l.Status,
		&l.CreatedAt,
		&l.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return l, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.lead.GetLeadByID r.db.Query", zap.Error(err))
		return l, err
	}

	r.logger.Info("pkg.repo.lead.GetLeadByID", zap.Any("lead", l))

	return l, nil
}

func (r *repo) UpdateLeadStatus(ctx context.Context, id int, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE lead SET status = $2, updated_at = now() WHERE id = $1;`, id, status)
	if err != nil {
		r.logger.Error("pkg.repo.lead.UpdateLeadStatus r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
