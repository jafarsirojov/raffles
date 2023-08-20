package lending

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (r *repo) SaveAvailability(ctx context.Context, data structs.Availability) (err error) {
	_, err = r.db.Exec(ctx, `
INSERT INTO availability(
    lending_id,
    price_aed,
    price_usd,
    unique_id,
    bedroom,
    parking,
    area,
    plot,
    special_gift) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);`,
		data.LendingID,
		data.Price.AED,
		data.Price.USD,
		data.UniqueID,
		data.Bedroom,
		data.Parking,
		data.Area,
		data.Plot,
		data.SpecialGift,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.SaveAvailability r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) GetAvailabilitiesByLandingID(ctx context.Context, landingID int) (list []structs.Availability, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    id,
    lending_id,
    price_aed,
    price_usd,
    unique_id,
    bedroom,
    parking,
    area,
    plot,
    special_gift
FROM availability
WHERE lending_id = $1;`, landingID)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.GetAvailabilitiesByLandingID r.db.Query", zap.Error(err))
		return list, err
	}

	for rows.Next() {
		var a structs.Availability
		err = rows.Scan(
			&a.ID,
			&a.LendingID,
			&a.Price.AED,
			&a.Price.USD,
			&a.UniqueID,
			&a.Bedroom,
			&a.Parking,
			&a.Area,
			&a.Plot,
			&a.SpecialGift,
		)
		if err != nil {
			r.logger.Error("pkg.repo.admin.lending.GetAvailabilitiesByLandingID rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, a)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	return list, nil
}

func (r *repo) UpdateAvailability(ctx context.Context, data structs.Availability) (err error) {
	_, err = r.db.Exec(ctx, `
UPDATE availability SET
    lending_id = $2,
    price_aed = $3,
    price_usd = $4,
    unique_id = $5,
    bedroom = $6,
    parking = $7,
    area = $8,
    plot = $9,
    special_gift = $10,
    updated_at = now()
    WHERE id = $1`,
		data.LendingID,
		data.Price.AED,
		data.Price.USD,
		data.UniqueID,
		data.Bedroom,
		data.Parking,
		data.Area,
		data.Plot,
		data.SpecialGift,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.UpdateAvailability r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) SelectPaymentPlanByAvailabilityID(ctx context.Context, id int) (paymentPlan string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT payment_plan FROM availability WHERE id = $1;`, id).Scan(&paymentPlan)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.SelectPaymentPlanByAvailability r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return paymentPlan, err
	}

	return paymentPlan, nil
}

func (r *repo) UpdatePaymentPlan(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE availability SET payment_plan = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.UpdatePaymentPlan r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}
