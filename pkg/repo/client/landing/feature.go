package landing

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (r *repo) SelectFeaturesAndAmenitiesByIDs(ctx context.Context, ids []int) (list []structs.FeatureOrAmenity, err error) {
	rows, err := r.db.Query(ctx, `
SELECT 
    name,
    icon
FROM feature_or_amenity
WHERE id = ANY ($1);`, ids)
	if err != nil {
		r.logger.Error("pkg.repo.client.landing.SelectFeaturesAndAmenitiesByIDs r.db.Query", zap.Error(err))
		return list, err
	}

	for rows.Next() {
		var data structs.FeatureOrAmenity
		err = rows.Scan(
			&data.Name,
			&data.Icon,
		)
		if err != nil {
			r.logger.Error("pkg.repo.client.landing.SelectFeaturesAndAmenitiesByIDs rows.Scan()", zap.Error(err))
			return nil, err
		}

		list = append(list, data)
	}

	if len(list) == 0 {
		return nil, errors.ErrNotFound
	}

	return list, nil
}
