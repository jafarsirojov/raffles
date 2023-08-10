package auth

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
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

func NewRepo(params Params) interfaces.AuthClientRepo {
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
			r.logger.Warn("pkg.repo.client.auth.SignIn r.db.QueryRow - not found", zap.String("login", login))
			return token, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.client.auth.SignIn r.db.QueryRow", zap.Error(err))
		return token, err
	}

	if password == passwordForDB {
		return token, nil
	}

	return token, errors.ErrInvalidPassword
}

func (r *repo) SaveUser(ctx context.Context, user structs.SignUp) (token string, err error) {

	_, err = r.db.Exec(ctx,
		`INSERT INTO user (first_name, last_name, phone, login, password, token) VALUES ($1, $2, $3, $4, $5);`,
		user.FirstName,
		user.LastName,
		user.Phone,
		user.Login,
		user.Password,
		user.Token,
	)
	if err != nil {
		r.logger.Error("pkg.repo.client.auth.SaveUser r.db.Exec", zap.Error(err))
		return token, err
	}

	return token, nil
}

func (r *repo) CheckToken(ctx context.Context, token string) (id int, login string, err error) {

	err = r.db.QueryRow(ctx, `
SELECT id, login
FROM user
WHERE status = 'enabled' AND token = $1 order by id;`, token).Scan(&id, &login)
	if err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Warn("pkg.repo.client.auth.CheckToken r.db.QueryRow - not found", zap.String("login", login))
			return id, login, errors.ErrNotFound
		}
		r.logger.Error("pkg.repo.client.auth.CheckToken r.db.QueryRow", zap.String("token", token), zap.Error(err))
		return id, login, err
	}

	return id, login, nil
}
