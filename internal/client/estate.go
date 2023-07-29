package client

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) GetEstates(ctx context.Context, offset, limit int, options structs.SearchOptionsDTO) (estates []structs.EstateForList, err error) {
	estates, err = s.estateRepo.SelectEstates(ctx, offset, limit)
	if err != nil {
		if err != errors.ErrNotFound {
			s.logger.Info("internal.client.GetEstates s.estateRepo.SelectEstates not found",
				zap.Int("offset", offset), zap.Int("limit", limit))
			return estates, err
		}

		s.logger.Error("internal.client.GetEstates s.estateRepo.SelectEstates",
			zap.Error(err), zap.Int("offset", offset), zap.Int("limit", limit))
		return estates, err
	}

	return estates, nil
}

func (s *service) GetLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, err error) {
	estates, err = s.estateRepo.SelectLuxuryEstates(ctx, offset, limit)
	if err != nil {
		if err != errors.ErrNotFound {
			s.logger.Info("internal.client.GetLuxuryEstates s.estateRepo.SelectLuxuryEstates not found",
				zap.Int("offset", offset), zap.Int("limit", limit))
			return estates, err
		}

		s.logger.Error("internal.client.GetLuxuryEstates s.estateRepo.SelectLuxuryEstates",
			zap.Error(err), zap.Int("offset", offset), zap.Int("limit", limit))
		return estates, err
	}

	return estates, nil
}

func (s *service) GetEstateByID(ctx context.Context, id int) (estate structs.Estate, err error) {
	estate, err = s.estateRepo.SelectEstateByID(ctx, id)
	if err != nil {
		if err != errors.ErrNotFound {
			s.logger.Info("internal.client.GetEstateByID s.estateRepo.SelectEstateByID not found",
				zap.Int("id", id))
			return estate, err
		}

		s.logger.Error("internal.client.GetEstateByID s.estateRepo.SelectEstateByID",
			zap.Error(err), zap.Int("id", id))
		return estate, err
	}

	return estate, nil
}

var Cache = make(map[string]structs.SearchOptions)

const cacheKey = "option"

func (s *service) GetSearchOptions(ctx context.Context) (options structs.SearchOptions, err error) {

	options, ok := Cache[cacheKey]
	if !ok {
		options, err = s.estateRepo.SelectSearchOptions(ctx)
		if err != nil {
			if err != errors.ErrNotFound {
				s.logger.Info("internal.client.GetSearchOptions s.estateRepo.SelectSearchOptions not found")
				return options, err
			}

			s.logger.Error("internal.client.GetSearchOptions s.estateRepo.SelectSearchOptions", zap.Error(err))
			return options, err
		}

		options.PropertyTypes = []structs.PropertyTypeKeyValue{
			{1, "property"},
			{2, "villa"},
			{3, "dacha"},
		}

		Cache[cacheKey] = options
	}

	return options, nil
}

func (s *service) ClearCache() {
	delete(Cache, cacheKey)
}
