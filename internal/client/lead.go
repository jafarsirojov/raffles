package client

import (
	"context"
	"crm/internal/structs"
	"go.uber.org/zap"
)

func (s *service) SaveLead(ctx context.Context, request structs.Lead) error {

	err := s.leadRepo.InsertLead(ctx, request)
	if err != nil {
		s.logger.Error("internal.client.SaveLead s.leadRepo.InsertLead", zap.Error(err))
		return err
	}

	return nil
}
