package favorites

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

func NewRepo(params Params) interfaces.FavoritesClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) InsertFavorite(ctx context.Context, favorite structs.Favorite) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO favorite (user_id, estate_id)VALUES ($1, $2)", favorite.UserID, favorite.EstateID)
	if err != nil {
		r.logger.Error("pkg.client.repo.favorites.InsertFavorite r.db.Exec",
			zap.Error(err), zap.Any("favorite", favorite))
		return err
	}

	return nil
}

func (r *repo) DeleteFavorite(ctx context.Context, favorite structs.Favorite) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM favorite WHERE user_id = $1 AND estate_id = $2;", favorite.UserID, favorite.EstateID)
	if err != nil {
		r.logger.Error("pkg.client.repo.favorites.DeleteFavorite r.db.Exec",
			zap.Error(err), zap.Any("favorite", favorite))
		return err
	}

	return nil
}

func (r *repo) GetFavorites(ctx context.Context, userID int) (favorites []int, err error) {
	rows, err := r.db.Query(ctx, "SELECT estate_id FROM favorite WHERE user_id = $1", userID)
	if err != nil {
		r.logger.Error("pkg.client.repo.favorites.GetFavorites r.db.Query",
			zap.Error(err), zap.Int("userID", userID))
		return nil, err
	}

	for rows.Next() {
		var f int
		err = rows.Scan(&f)
		if err != nil {
			r.logger.Error("pkg.client.repo.favorites.GetFavorites rows.Scan",
				zap.Error(err), zap.Int("userID", userID))
			return nil, err
		}

		favorites = append(favorites, f)
	}
	if len(favorites) == 0 {
		return nil, errors.ErrNotFound
	}

	return favorites, nil
}
