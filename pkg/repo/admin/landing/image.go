package landing

import (
	"context"
	"go.uber.org/zap"
)

func (r *repo) GetImagesByLandingID(ctx context.Context, id int) (images []string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT images FROM landing WHERE id = $1;`, id).Scan(&images)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.GetImagesByLandingID r.db.QueryRow", zap.Error(err))
		return nil, err
	}

	return images, nil
}

func (r *repo) UpdateLandingImages(ctx context.Context, id int, images []string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET images = $2, updated_at = now() WHERE id = $1;`, id, images)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateLandingImages r.db.Exec", zap.Error(err))
		return err
	}

	return nil
}

func (r *repo) SelectBackgroundImageByLandingID(ctx context.Context, id int) (backgroundImage string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT background_image FROM landing WHERE id = $1;`, id).Scan(&backgroundImage)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectBackgroundImageByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return backgroundImage, err
	}

	return backgroundImage, nil
}

func (r *repo) UpdateBackgroundImage(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET background_image = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateBackgroundImage r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}

func (r *repo) SelectVideoCoverByLandingID(ctx context.Context, id int) (cover string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT video_cover FROM landing WHERE id = $1;`, id).Scan(&cover)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectVideoCoverByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return cover, err
	}

	return cover, nil
}

func (r *repo) UpdateVideoCover(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET video_cover = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateVideoCover r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}

func (r *repo) SelectBackgroundForMobileByLandingID(ctx context.Context, id int) (backgroundImage string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT background_for_mobile FROM landing WHERE id = $1;`, id).Scan(&backgroundImage)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectBackgroundForMobileByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return backgroundImage, err
	}

	return backgroundImage, nil
}

func (r *repo) UpdateBackgroundForMobile(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET background_for_mobile = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateBackgroundForMobile r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}

func (r *repo) SelectMainLogoByLandingID(ctx context.Context, id int) (logo string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT main_logo FROM landing WHERE id = $1;`, id).Scan(&logo)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectMainLogoByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return logo, err
	}

	return logo, nil
}

func (r *repo) UpdateMainLogo(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET main_logo = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateMainLogo r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}

func (r *repo) SelectPartnerLogoByLandingID(ctx context.Context, id int) (logo string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT partner_logo FROM landing WHERE id = $1;`, id).Scan(&logo)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectPartnerLogoLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return logo, err
	}

	return logo, nil
}

func (r *repo) UpdatePartnerLogo(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET partner_logo = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdatePartnerLogo r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}

func (r *repo) SelectOurLogoByLandingID(ctx context.Context, id int) (logo string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT our_logo FROM landing WHERE id = $1;`, id).Scan(&logo)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.SelectOurLogoByLandingID r.db.QueryRow",
			zap.Error(err), zap.Int("id", id))
		return logo, err
	}

	return logo, nil
}

func (r *repo) UpdateOurLogo(ctx context.Context, id int, new string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE landing SET our_logo = $2, updated_at = now() WHERE id = $1;`, id, new)
	if err != nil {
		r.logger.Error("pkg.repo.admin.landing.UpdateOurLogo r.db.Exec",
			zap.Error(err), zap.Int("id", id), zap.Any("new", new))
		return err
	}

	return nil
}
