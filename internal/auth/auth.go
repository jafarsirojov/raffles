package auth

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewService)

type Params struct {
	fx.In
	Logger   *zap.Logger
	AuthRepo interfaces.AuthRepo
}

type Service interface {
	SignIn(ctx context.Context, in structs.SignIn) (token string, err error)
	CheckAdminAuthToken(ctx context.Context, token string) (id int, login string, err error)
}

type service struct {
	logger   *zap.Logger
	authRepo interfaces.AuthRepo
}

func NewService(params Params) Service {
	return &service{
		logger:   params.Logger,
		authRepo: params.AuthRepo,
	}
}

func (s *service) SignIn(ctx context.Context, in structs.SignIn) (token string, err error) {
	token, err = s.authRepo.SignIn(ctx, in.Login, in.Password)
	if err != nil {
		if err == errors.ErrNotFound || err == errors.ErrInvalidPassword {
			s.logger.Warn("internal.auth.SignIn s.authRepo.SignIn", zap.Error(err), zap.Any("login", in.Login))
			return token, err
		}

		s.logger.Error("internal.auth.SignIn s.authRepo.SignIn", zap.Error(err))
		return token, err
	}

	s.logger.Info("internal.auth.SignIn - successfully", zap.Any("login", in.Login))

	return token, nil
}

func (s *service) CheckAdminAuthToken(ctx context.Context, token string) (id int, login string, err error) {
	id, login, err = s.authRepo.CheckAdminToken(ctx, token)
	if err != nil {
		if err == errors.ErrNotFound {
			s.logger.Warn("internal.auth.CheckAdminAuthToken s.authRepo.CheckAdminToken", zap.Error(err))
			return id, login, err
		}

		s.logger.Error("internal.auth.CheckAdminAuthToken s.authRepo.CheckAdminToken", zap.Error(err))
		return id, login, err
	}

	return id, login, nil
}
