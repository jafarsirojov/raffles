package interfaces

import (
	"context"
	"crm/internal/structs"
)

type FavoritesClientRepo interface {
	InsertFavorite(ctx context.Context, favorite structs.Favorite) error
	DeleteFavorite(ctx context.Context, favorite structs.Favorite) error
	GetFavorites(ctx context.Context, userID int) (favorites []int, err error)
}
