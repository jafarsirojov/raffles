package interfaces

import (
	"context"
)

type AuthRepo interface {
	SignIn(ctx context.Context, login, password string) (token string, err error)
	CheckAdminToken(ctx context.Context, token string) (id int, login string, err error)
}
