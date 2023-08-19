package admin

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
)

func (s *service) GetLendingList(ctx context.Context) (list []structs.LendingList, err error) {
	list, err = s.lendingRepo.GetLendingList(ctx)
	if err != nil {
		s.logger.Error("internal.admin.GetLendingList s.lendingRepo.GetLendingList", zap.Error(err))
		return nil, err
	}

	return list, nil
}

func (s *service) GetLendingData(ctx context.Context, landingID int) (data structs.Lending, err error) {
	data, err = s.lendingRepo.GetLendingByID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.GetLendingData s.lendingRepo.GetLendingByID",
			zap.Error(err), zap.Int("landingID", landingID))
		return data, err
	}

	data.FeaturesAndAmenities, err = s.lendingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, data.FeaturesAndAmenitiesIDs)
	if err != nil {
		s.logger.Error("internal.admin.GetLendingData s.lendingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Int("landingID", landingID))
		return data, err
	}

	return data, nil
}

func (s *service) SaveLending(ctx context.Context, data structs.Lending) error {
	err := s.lendingRepo.SaveLending(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.SaveLending s.lendingRepo.SaveLending",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) UpdateLending(ctx context.Context, data structs.Lending) error {
	err := s.lendingRepo.UpdateLending(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.UpdateLending s.lendingRepo.UpdateLending",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) UploadLendingImages(ctx context.Context, id int, files []multipart.File) error {

	newImagesName, err := s.uploadImages(files)
	if err != nil {
		s.logger.Error("internal.admin.UploadLendingImages s.uploadImages", zap.Error(err), zap.Int("id", id))
		return err
	}

	imagesName, err := s.lendingRepo.GetImagesByLendingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.UploadLendingImages s.lendingRepo.GetImagesByLendingID",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	imagesName = append(imagesName, newImagesName...)

	err = s.lendingRepo.UpdateLendingImages(ctx, id, imagesName)
	if err != nil {
		s.logger.Error("internal.admin.UploadLendingImages s.lendingRepo.UpdateLendingImages",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *service) DeleteLendingImages(ctx context.Context, id int, imageName string) error {
	imageNames, err := s.lendingRepo.GetImagesByLendingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLendingImages s.lendingRepo.GetImagesByLendingID",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	exist := false
	index := 0
	for i, name := range imageNames {
		if name == imageName {
			index = i
			exist = true
			break
		}
	}

	if exist {
		imageNames = append(imageNames[:index], imageNames[index+1:]...)
	} else {
		return errors.ErrNotFound
	}

	err = s.lendingRepo.UpdateLendingImages(ctx, id, imageNames)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLendingImages s.lendingRepo.UpdateLendingImages",
			zap.Error(err), zap.Int("id", id), zap.Any("images", imageNames))
		return err
	}

	err = os.Remove(structs.ImagePath + imageName)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLendingImages os.Remove",
			zap.Error(err), zap.Int("id", id), zap.Any("image", imageName))
		return err
	}

	return nil
}

func (s *service) GetFeaturesAndAmenities(ctx context.Context) (list []structs.FeatureOrAmenity, err error) {
	list, err = s.lendingRepo.SelectFeaturesAndAmenities(ctx)
	if err != nil {
		s.logger.Error("internal.admin.GetFeaturesAndAmenities s.lendingRepo.SelectFeaturesAndAmenities",
			zap.Error(err))
		return nil, err
	}

	return list, nil
}
