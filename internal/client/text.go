package client

import (
	"context"
	"crm/internal/structs"
	"go.uber.org/zap"
)

func (s *service) GetTexts(ctx context.Context) (texts []structs.Text, err error) {
	texts, err = s.textRepo.GetTexts(ctx)
	if err != nil {
		s.logger.Error("internal.client.GetTexts s.textRepo.GetTexts", zap.Error(err))
		return nil, err
	}

	return nil, err
}
