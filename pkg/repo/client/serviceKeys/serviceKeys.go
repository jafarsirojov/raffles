package serviceKeys

import (
	"context"
	interfaces "crm/internal/interface"
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

func NewRepo(params Params) interfaces.ServiceKeysClientRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SelectLendingIdByKey(ctx context.Context, key string) (id int, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT id FROM service_keys WHERE key = $1;`, key).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("pkg.repo.client.serviceKeys.SelectLendingIdByKey r.db.QueryRow not found",
				zap.String("key", key))
			return id, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.client.serviceKeys.SelectLendingIdByKey r.db.QueryRow",
			zap.String("key", key), zap.Error(err))
		return id, err
	}

	return id, nil
}
