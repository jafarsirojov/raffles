package interfaces

import "context"

type ServiceKeysClientRepo interface {
	SelectLandingIdByKey(ctx context.Context, key string) (id int, err error)
}
