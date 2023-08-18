package interfaces

import (
	"context"
	"crm/internal/structs"
)

type LendingAdminRepo interface {
	SaveLending(ctx context.Context, data structs.Lending) error
	UpdateLending(ctx context.Context, data structs.Lending) (err error)
	GetLendingList(ctx context.Context) (list []structs.LendingList, err error)
	GetLendingByID(ctx context.Context, id int) (data structs.Lending, err error)
}

type LendingClientRepo interface {
}
