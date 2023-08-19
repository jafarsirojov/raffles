package admin

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"go.uber.org/zap"
	"mime/multipart"
	"os"
)

func (s *service) GetLendingData(ctx context.Context, landingID string) (data structs.Lending, err error) {

	return data, nil
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
