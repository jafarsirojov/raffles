package client

import (
	"context"
	"crm/internal/structs"
	"go.uber.org/zap"
)

func (s *service) SaveFavorite(ctx context.Context, favorite structs.Favorite) error {
	err := s.favoritesRepo.InsertFavorite(ctx, favorite)
	if err != nil {
		s.logger.Error("internal.client.SaveFavorite s.favoritesRepo.InsertFavorite",
			zap.Error(err), zap.Any("favorite", favorite))
		return err
	}

	return nil
}

func (s *service) DeleteFavorite(ctx context.Context, favorite structs.Favorite) error {
	err := s.favoritesRepo.DeleteFavorite(ctx, favorite)
	if err != nil {
		s.logger.Error("internal.client.DeleteFavorite s.favoritesRepo.DeleteFavorite",
			zap.Error(err), zap.Any("favorite", favorite))
		return err
	}

	return nil
}

func (s *service) GetEstateFavorites(ctx context.Context, userID int) (estates []structs.EstateForList, err error) {
	favorites, err := s.favoritesRepo.GetFavorites(ctx, userID)
	if err != nil {
		s.logger.Error("internal.client.GetEstateFavorites s.favoritesRepo.GetFavorites",
			zap.Error(err), zap.Int("userID", userID))
		return nil, err
	}

	estates, err = s.estateRepo.SelectFavoriteEstates(ctx, favorites)
	if err != nil {
		s.logger.Error("internal.client.GetEstateFavorites s.estateRepo.SelectFavoriteEstates",
			zap.Error(err), zap.Any("favorites", favorites))
		return nil, err
	}

	return estates, nil
}
