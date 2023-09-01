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
    special_gift,
    special_gift_icon) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		data.LendingID,
		data.Price.AED,
		data.Price.USD,
		data.UniqueID,
		data.Bedroom,
		data.Parking,
		data.Area,
		data.Plot,
		data.SpecialGift,
		data.SpecialGiftIcon,
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
    special_gift,
    special_gift_icon
FROM availability
WHERE lending_id = $1
ORDER BY id;`, landingID)
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
			&a.SpecialGiftIcon,
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
    special_gift_icon = $11,
    updated_at = now()
    WHERE id = $1`,
		data.ID,
		data.LendingID,
		data.Price.AED,
		data.Price.USD,
		data.UniqueID,
		data.Bedroom,
		data.Parking,
		data.Area,
		data.Plot,
		data.SpecialGift,
		data.SpecialGiftIcon,
	)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.UpdateAvailability r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) DeleteAvailability(ctx context.Context, id int) (err error) {
	_, err = r.db.Exec(ctx, `DELETE FROM availability WHERE id = $1;`, id)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.DeleteAvailability r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
