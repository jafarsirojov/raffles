package interfaces

import (
	"context"
	"crm/internal/structs"
)

type EstateAdminRepo interface {
	GetEstateByID(ctx context.Context, id int) (structs.Estate, error)
	GetEstates(ctx context.Context, offset, limit int) ([]structs.Estate, error)
	GetEstatesByStatus(ctx context.Context, offset, limit int, status string) ([]structs.Estate, error)
	SaveEstate(ctx context.Context, request structs.Estate) error
	UpdateEstate(ctx context.Context, request structs.Estate) error
	UpdateEstateStatus(ctx context.Context, id int, status structs.Status) error
	GetImagesByEstateID(ctx context.Context, id int) ([]string, error)
	UpdateEstateImages(ctx context.Context, id int, images []string) error
}

type EstateClientRepo interface {
	SelectEstates(ctx context.Context, offset, limit int, options structs.SearchOptionsDTO) (estates []structs.EstateForList, err error)
	SelectLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error)
	SelectEstateByID(ctx context.Context, id int) (estate structs.Estate, err error)
	SelectSearchOptions(ctx context.Context) (options structs.SearchOptions, err error)
}
