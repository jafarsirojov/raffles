package client

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) GetEstates(ctx context.Context, offset, limit int, options structs.SearchOptionsDTO) (estates []structs.EstateForList, totalCount int, err error) {
	estates, err = s.estateRepo.SelectEstates(ctx, offset, limit, options)
	if err != nil {
		if err != errors.ErrNotFound {
			s.logger.Info("internal.client.GetEstates s.estateRepo.SelectEstates not found",
				zap.Int("offset", offset), zap.Int("limit", limit))
			return estates, totalCount, err
		}

		s.logger.Error("internal.client.GetEstates s.estateRepo.SelectEstates",
			zap.Error(err), zap.Int("offset", offset), zap.Int("limit", limit))
		return estates, totalCount, err
	}

	totalCount, err = s.estateRepo.GetEstatesTotalCount(ctx, false)
	if err != nil {
		s.logger.Error("internal.client.GetEstates s.estateRepo.GetEstatesTotalCount", zap.Error(err))
		return estates, totalCount, err
	}

	return estates, totalCount, err
}

func (s *service) GetLuxuryEstates(ctx context.Context, offset, limit int) (estates []structs.EstateForList, totalCount int, err error) {
	estates, err = s.estateRepo.SelectLuxuryEstates(ctx, offset, limit)
	if err != nil {
		if err != errors.ErrNotFound {
			s.logger.Info("internal.client.GetLuxuryEstates s.estateRepo.SelectLuxuryEstates not found",
				zap.Int("offset", offset), zap.Int("limit", limit))
			return estates, totalCount, err
		}

		s.logger.Error("internal.client.GetLuxuryEstates s.estateRepo.SelectLuxuryEstates",
			zap.Error(err), zap.Int("offset", offset), zap.Int("limit", limit))
		return estates, totalCount, err
	}

	totalCount, err = s.estateRepo.GetEstatesTotalCount(ctx, true)
	if err != nil {
		s.logger.Error("internal.client.GetEstates s.estateRepo.GetEstatesTotalCount", zap.Error(err))
		return estates, totalCount, err
	}

	return estates, totalCount, nil
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
			{1, "Co-op"},
			{2, "Condo"},
			{3, "Commercial"},
			{4, "Land"},
			{5, "Multi-family"},
			{6, "Manufactured"},
			{7, "Residential"},
			{8, "Townhouse"},
			{9, "Others"},
		}

		options.Status = []structs.Status{
			"active",
			"soled",
			"canceled",
			"new",
		}

		Cache[cacheKey] = options
	}

	return options, nil
}

func (s *service) ClearCache() {
	delete(Cache, cacheKey)
}

func (s *service) GetImageBaseURL() string {
	return structs.ImageBaseURL
}
