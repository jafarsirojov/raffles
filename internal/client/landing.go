package client

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
)

func (s *service) GetLandingsList(ctx context.Context) (list structs.ListMainPage, err error) {

	list.List, err = s.landingRepo.SelectLandingList(ctx)
	if err != nil {
		s.logger.Error("internal.client.GetLandingsList s.landingRepo.SelectLandingList",
			zap.Error(err))
		return list, err
	}

	list.FileURL = structs.FileBaseURLPathRafflesHomes

	return list, err
}

func (s *service) GetLandingData(ctx context.Context, key string) (data structs.LandingData, err error) {

	id, err := s.serviceKeyRepo.SelectLandingIdByKey(ctx, key)
	if err != nil {
		if err == errors.ErrNotFound {
			s.logger.Warn("internal.client.GetLandingData s.serviceKeyRepo.SelectLandingIdByKey not found",
				zap.String("key", key))
			return data, errors.ErrBadRequest
		}
		s.logger.Error("internal.client.GetLandingData s.serviceKeyRepo.SelectLandingIdByKey",
			zap.Error(err))
		return data, err
	}

	landing, err := s.landingRepo.SelectLandingData(ctx, id)
	if err != nil {
		s.logger.Error("internal.client.GetLandingData s.serviceKeyRepo.SelectLandingData",
			zap.Error(err))
		return data, err
	}

	data.Name = landing.Name
	data.MainDescription = landing.MainDescription
	data.FullName = landing.FullName
	data.Slogan = landing.Slogan
	data.Address = landing.Address
	data.StartingPrice = landing.StartingPrice
	data.ListingDetails = landing.ListingDetails
	data.FeaturesAndAmenities = landing.FeaturesAndAmenities
	data.Title = landing.Title
	data.Description = landing.Description
	data.Video = landing.Video
	data.VideoCover = landing.VideoCover
	data.FilePlan = landing.FilePlan
	data.TitlePlan = landing.TitlePlan
	data.Images = landing.Images
	data.BackgroundImage = landing.BackgroundImage
	data.BackgroundForMobile = landing.BackgroundForMobile
	data.MainLogo = landing.MainLogo
	data.PartnerLogo = landing.PartnerLogo
	data.OurLogo = landing.OurLogo
	data.Latitude = landing.Latitude
	data.Longitude = landing.Longitude
	data.LocationDescription = landing.LocationDescription

	data.FeaturesAndAmenities, err = s.landingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, landing.FeaturesAndAmenitiesIDs)
	if err != nil {
		s.logger.Error("internal.client.GetLandingData s.landingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Any("landing.FeaturesAndAmenitiesIDs", landing.FeaturesAndAmenitiesIDs))
	}

	data.Availabilities, err = s.landingRepo.GetAvailabilitiesByLandingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.client.GetLandingData s.landingRepo.GetAvailabilitiesByLandingID",
			zap.Error(err), zap.Any("landing.FeaturesAndAmenitiesIDs", landing.FeaturesAndAmenitiesIDs))
	}

	data.FileURL = structs.FileBaseURLPathRafflesHomes

	return data, nil
}
