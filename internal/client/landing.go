package client

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) GetLandingsList(ctx context.Context) (data structs.LendingListMainPage, err error) {

	data, err = s.lendingRepo.SelectLandingList(ctx)
	if err != nil {
		s.logger.Error("internal.client.GetLendingData s.lendingRepo.SelectLendingData",
			zap.Error(err))
		return data, err
	}

	data.FileURL = structs.FileBaseURLPathRafflesHomes

	return data, err
}

func (s *service) GetLendingData(ctx context.Context, key string) (data structs.LendingData, err error) {

	id, err := s.serviceKeyRepo.SelectLendingIdByKey(ctx, key)
	if err != nil {
		if err == errors.ErrNotFound {
			s.logger.Warn("internal.client.GetLendingData s.lendingRepo.SelectLendingIdByKey not found",
				zap.String("key", key))
			return data, errors.ErrBadRequest
		}
		s.logger.Error("internal.client.GetLendingData s.lendingRepo.SelectLendingIdByKey",
			zap.Error(err))
		return data, err
	}

	lending, err := s.lendingRepo.SelectLendingData(ctx, id)
	if err != nil {
		s.logger.Error("internal.client.GetLendingData s.lendingRepo.SelectLendingData",
			zap.Error(err))
		return data, err
	}

	data.Name = lending.Name
	data.MainDescription = lending.MainDescription
	data.FullName = lending.FullName
	data.Address = lending.Address
	data.StartingPrice = lending.StartingPrice
	data.ListingDetails = lending.ListingDetails
	data.FeaturesAndAmenities = lending.FeaturesAndAmenities
	data.Title = lending.Title
	data.Description = lending.Description
	data.Video = lending.Video
	data.Images = lending.Images
	data.BackgroundImage = lending.BackgroundImage

	data.FeaturesAndAmenities, err = s.lendingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, lending.FeaturesAndAmenitiesIDs)
	if err != nil {
		s.logger.Error("internal.client.GetLendingData s.lendingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Any("lending.FeaturesAndAmenitiesIDs", lending.FeaturesAndAmenitiesIDs))
		return data, err
	}

	data.Availabilities, err = s.lendingRepo.GetAvailabilitiesByLandingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.client.GetLendingData s.lendingRepo.GetAvailabilitiesByLandingID",
			zap.Error(err), zap.Any("lending.FeaturesAndAmenitiesIDs", lending.FeaturesAndAmenitiesIDs))
		return data, err
	}

	data.FileURL = structs.FileBaseURLPathRafflesHomes

	return data, err
}
