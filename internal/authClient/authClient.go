package authClient

import (
	"context"
	interfaces "crm/internal/interface"
	"crm/internal/structs"
	"crm/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(NewService)

type Params struct {
	fx.In
	Logger   *zap.Logger
	AuthRepo interfaces.AuthClientRepo
}

type Service interface {
	SignIn(ctx context.Context, in structs.SignIn) (token string, err error)
	SignUp(ctx context.Context, up structs.SignUp) (token string, err error)
	CheckAuthToken(ctx context.Context, token string) (id int, login string, err error)
}

type service struct {
	logger   *zap.Logger
	authRepo interfaces.AuthClientRepo
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
			s.logger.Warn("internal.authClient.SignIn s.authRepo.SignIn", zap.Error(err), zap.Any("login", in.Login))
			return token, err
		}

		s.logger.Error("internal.authClient.SignIn s.authRepo.SignIn", zap.Error(err))
		return token, err
	}

	s.logger.Info("internal.authClient.SignIn - successfully", zap.Any("login", in.Login))

	return token, nil
}

func (s *service) SignUp(ctx context.Context, up structs.SignUp) (token string, err error) {

	up.Token = uuid.NewString()

	token, err = s.authRepo.SaveUser(ctx, up)
	if err != nil {
		s.logger.Error("internal.authClient.SignUp s.authRepo.SaveUser", zap.Error(err))
		return token, err
	}

	s.logger.Info("internal.authClient.SignUp - successfully", zap.Any("login", up.Login))

	return token, nil
}

func (s *service) CheckAuthToken(ctx context.Context, token string) (id int, login string, err error) {
	id, login, err = s.authRepo.CheckToken(ctx, token)
	if err != nil {
		if err == errors.ErrNotFound {
			s.logger.Warn("internal.authClient.CheckAuthToken s.authRepo.CheckToken", zap.Error(err))
			return id, login, err
		}

		s.logger.Error("internal.authClient.CheckAuthToken s.authRepo.CheckToken", zap.Error(err))
		return id, login, err
	}

	return id, login, nil
}
