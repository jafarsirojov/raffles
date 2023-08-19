package interfaces

import "context"

type ServiceKeysClientRepo interface {
	SelectLendingIdByKey(ctx context.Context, key string) (id int, err error)
}
