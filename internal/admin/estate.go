package admin

import (
	"context"
	"crm/internal/structs"
	"crm/pkg/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"os"
	"strings"
)

func (s *service) GetEstateByID(ctx context.Context, id int) (estate structs.Estate, err error) {
	estate, err = s.estateRepo.GetEstateByID(ctx, id)
	if err != nil {
		if err == errors.ErrNotFound {
			s.logger.Warn("internal.admin.GetEstateByID s.estateRepo.GetEstateByID not found", zap.Int("id", id))
			return structs.Estate{}, err
		}
		s.logger.Error("internal.admin.GetEstateByID s.estateRepo.GetEstateByID",
			zap.Int("id", id), zap.Error(err))
		return structs.Estate{}, err
	}

	return estate, nil
}

func (s *service) GetEstates(ctx context.Context, offset, limit int, status string) (estates []structs.EstateForList, totalCount int, err error) {

	if len(strings.TrimSpace(status)) == 0 {
		estates, err = s.estateRepo.GetEstatesForList(ctx, offset, limit)
		if err != nil {
			if err == errors.ErrNotFound {
				s.logger.Warn("internal.admin.GetEstates s.estateRepo.GetEstatesForList not found",
					zap.Int("offset", offset), zap.Int("limit", limit))
				return nil, totalCount, err
			}
			s.logger.Error("internal.admin.GetEstates s.estateRepo.GetEstatesForList", zap.Error(err),
				zap.Int("offset", offset), zap.Int("limit", limit))
			return nil, totalCount, err
		}
	} else {
		estates, err = s.estateRepo.GetEstatesForListByStatus(ctx, offset, limit, status)
		if err != nil {
			if err == errors.ErrNotFound {
				s.logger.Warn("internal.admin.GetEstates s.estateRepo.GetEstatesForListByStatus not found",
					zap.String("status", status), zap.Int("offset", offset), zap.Int("limit", limit))
				return nil, totalCount, err
			}
			s.logger.Error("internal.admin.GetEstates s.estateRepo.GetEstatesForListByStatus", zap.Error(err),
				zap.String("status", status), zap.Int("offset", offset), zap.Int("limit", limit))
			return nil, totalCount, err
		}
	}

	totalCount, err = s.estateRepo.GetEstatesTotalCount(ctx, status)
	if err != nil {
		s.logger.Error("internal.admin.GetEstates s.estateRepo.GetEstatesTotalCount", zap.Error(err),
			zap.String("status", status))
		return estates, totalCount, err
	}

	return estates, totalCount, nil
}

func (s *service) AddEstate(ctx context.Context, request structs.Estate) error {
	err := s.estateRepo.SaveEstate(ctx, request)
	if err != nil {
		s.logger.Error("internal.admin.AddEstate s.estateRepo.SaveEstate",
			zap.Error(err), zap.Any("request", request))
		return err
	}

	return nil
}

func (s *service) UpdateEstate(ctx context.Context, request structs.Estate) error {
	err := s.estateRepo.UpdateEstate(ctx, request)
	if err != nil {
		s.logger.Error("internal.admin.UpdateEstate s.estateRepo.UpdateEstate",
			zap.Error(err), zap.Any("request", request))
		return err
	}

	return nil
}

func (s *service) DeleteEstate(ctx context.Context, id int) error {
	err := s.estateRepo.UpdateEstateStatus(ctx, id, "deleted")
	if err != nil {
		s.logger.Error("internal.admin.DeleteEstate s.estateRepo.UpdateEstateStatus",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *service) ApprovedEstate(ctx context.Context, id int) error {
	err := s.estateRepo.UpdateEstateStatus(ctx, id, "active")
	if err != nil {
		s.logger.Error("internal.admin.ApprovedEstate s.estateRepo.UpdateEstateStatus",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *service) UploadImages(ctx context.Context, id int, file *multipart.File) error {

	newUUID := uuid.NewString()
	filename := newUUID + ".png"

	newFile, err := os.Create(structs.ImagePath + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadImages os.Create", zap.Error(err), zap.Int("id", id))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(*file)
	if err != nil {
		s.logger.Error("internal.admin.UploadImages io.ReadAll", zap.Error(err), zap.Int("id", id))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadImages io.ReadAll", zap.Error(err), zap.Int("id", id))
		return err
	}

	imageNames, err := s.estateRepo.GetImagesByEstateID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.UploadImages s.estateRepo.GetImagesByEstateID",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	if len(imageNames) <= 1 && len(imageNames[0]) == 0 {
		imageNames = nil
	}

	imageNames = append(imageNames, filename)

	err = s.estateRepo.UpdateEstateImages(ctx, id, imageNames)
	if err != nil {
		s.logger.Error("internal.admin.UploadImages s.estateRepo.UpdateEstateImages",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *service) DeleteImages(ctx context.Context, id int, imageName string) error {
	imageNames, err := s.estateRepo.GetImagesByEstateID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.DeleteImages s.estateRepo.GetImagesByEstateID",
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

	err = s.estateRepo.UpdateEstateImages(ctx, id, imageNames)
	if err != nil {
		s.logger.Error("internal.admin.DeleteImages s.estateRepo.UpdateEstateImages",
			zap.Error(err), zap.Int("id", id), zap.Any("images", imageNames))
		return err
	}

	err = os.Remove(structs.ImagePath + imageName)
	if err != nil {
		s.logger.Error("internal.admin.DeleteImages os.Remove",
			zap.Error(err), zap.Int("id", id), zap.Any("image", imageName))
		return err
	}

	return nil
}

func (s *service) GetImageBaseURL() string {
	return structs.ImageBaseURL
}
