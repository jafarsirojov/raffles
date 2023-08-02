package admin

import (
	"context"
	"crm/internal/structs"
	"go.uber.org/zap"
)

func (s *service) GetTexts(ctx context.Context) (texts []structs.Text, err error) {
	texts, err = s.textRepo.GetTexts(ctx)
	if err != nil {
		s.logger.Error("internal.admin.GetTexts s.textRepo.GetTexts", zap.Error(err))
		return nil, err
	}

	return texts, err
}

func (s *service) UpdateText(ctx context.Context, text structs.Text) error {
	err := s.textRepo.UpdateText(ctx, text)
	if err != nil {
		s.logger.Error("internal.admin.UpdateText s.textRepo.UpdateText", zap.Error(err), zap.Any("text", text))
		return err
	}

	return nil
}
