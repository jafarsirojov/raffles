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

func (r *repo) SelectBackgroundImageLandingID(ctx context.Context, id int) (backgroundImage string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT background_image FROM lending WHERE id = $1;`, id).Scan(&backgroundImage)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.SelectBackgroundImageLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return backgroundImage, err
	}

	return backgroundImage, nil
}

func (r *repo) UpdateBackgroundImage(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE lending SET background_image = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.lending.UpdateBackgroundImage r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}
