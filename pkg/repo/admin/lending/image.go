package lending

import (
	"context"
	"go.uber.org/zap"
)

func (r *repo) GetImagesByLendingID(ctx context.Context, id int) (images []string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT images FROM lending WHERE id = $1;`, id).Scan(&images)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.GetImagesByLendingID r.db.QueryRow", zap.Error(err))
		return nil, err
	}

	return images, nil
}

func (r *repo) UpdateLendingImages(ctx context.Context, id int, images []string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE lending SET images = $2, updated_at = now() WHERE id = $1;`, id, images)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.UpdateLendingImages r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}
