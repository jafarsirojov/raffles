package interfaces

import (
	"context"
	"crm/internal/structs"
)

type AuthAdminRepo interface {
	SignIn(ctx context.Context, login, password string) (token string, err error)
	CheckAdminToken(ctx context.Context, token string) (id int, login string, err error)
}

type AuthClientRepo interface {
	SignIn(ctx context.Context, login, password string) (token string, err error)
	SaveUser(ctx context.Context, user structs.SignUp) (token string, err error)
	CheckToken(ctx context.Context, token string) (id int, login string, err error)
}
