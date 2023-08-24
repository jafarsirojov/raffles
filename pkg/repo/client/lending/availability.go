package lending

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

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
		r.logger.Error("pkg.repo.client.lending.GetAvailabilitiesByLandingID r.db.Query", zap.Error(err))
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
			r.logger.Error("pkg.repo.client.lending.GetAvailabilitiesByLandingID rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, a)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	return list, nil
}
