package auth

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

func NewRepo(params Params) interfaces.AuthRepo {
	return &repo{
		db:     params.DB,
		logger: params.Logger,
	}
}

type repo struct {
	db     db.Querier
	logger *zap.Logger
}

func (r *repo) SignIn(ctx context.Context, login, password string) (token string, err error) {
	var passwordForDB string
	err = r.db.QueryRow(ctx, `
SELECT password, token
FROM admin
WHERE status = 'enabled' AND login = $1;`, login).Scan(&passwordForDB, &token)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("pkg.repo.auth.SignIn r.db.QueryRow - not found", zap.String("login", login))
			return token, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.auth.SignIn r.db.QueryRow", zap.Error(err))
		return token, err
	}

	if password == passwordForDB {
		return token, nil
	}

	return token, errors.ErrInvalidPassword
}

func (r *repo) CheckAdminToken(ctx context.Context, token string) (id int, login string, err error) {

	err = r.db.QueryRow(ctx, `
SELECT id, login
FROM admin
WHERE status = 'enabled' AND token = $1 order by id;`, token).Scan(&id, &login)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("pkg.repo.auth.CheckAdminToken r.db.QueryRow - not found", zap.String("login", login))
			return id, login, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.auth.CheckAdminToken r.db.QueryRow", zap.String("token", token), zap.Error(err))
		return id, login, err
	}

	return id, login, nil
}
