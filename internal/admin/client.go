package admin

import (
	"context"
	"crm/internal/structs"
)

func (s *service) GetClientsFavorites(ctx context.Context, offset, limit int) (clientsFavorites []structs.Client, err error) {
	return clientsFavorites, nil
}
