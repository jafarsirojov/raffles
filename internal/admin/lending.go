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
	}

	data.Availabilities, err = s.lendingRepo.GetAvailabilitiesByLandingID(ctx, data.ID)
	if err != nil {
		s.logger.Error("internal.admin.GetLendingData s.lendingRepo.GetAvailabilitiesByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
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

func (s *service) SaveAvailability(ctx context.Context, data structs.Availability) error {
	err := s.lendingRepo.SaveAvailability(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.SaveAvailability s.lendingRepo.SaveAvailability",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) UpdateAvailability(ctx context.Context, data structs.Availability) error {
	err := s.lendingRepo.UpdateAvailability(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.UpdateAvailability s.lendingRepo.UpdateAvailability",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) RemoveAvailability(ctx context.Context, id int) error {
	err := s.lendingRepo.DeleteAvailability(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.RemoveAvailability s.lendingRepo.RemoveAvailability",
			zap.Error(err), zap.Any("id", id))
		return err
	}

	return nil
}

func (s *service) UploadLendingImages(ctx context.Context, id int, files []multipart.File) error {

	newImagesName, err := s.uploadImages(files, structs.FilePathRafflesHomes)
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

	err = os.Remove(structs.FilePathRafflesHomes + imageName)
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

func (s *service) UploadFilePlan(ctx context.Context, availabilityID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan io.ReadAll", zap.Error(err))
		return err
	}

	paymentPlan, err := s.lendingRepo.SelectFilePlanByLandingID(ctx, availabilityID)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan s.lendingRepo.SelectFilePlanByLandingID",
			zap.Error(err), zap.Int("availabilityID", availabilityID))
		return err
	}

	err = s.lendingRepo.UpdateFilePlan(ctx, availabilityID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan s.lendingRepo.UpdateFilePlan",
			zap.Error(err), zap.Int("availabilityID", availabilityID))
		return err
	}

	if len(strings.TrimSpace(paymentPlan)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + paymentPlan)
		if err != nil {
			s.logger.Error("internal.admin.UploadFilePlan os.Remove",
				zap.Error(err), zap.Int("availabilityID", availabilityID), zap.Any("old filename", paymentPlan))
		}
	}

	return nil
}

func (s *service) UploadBackgroundImage(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage io.ReadAll", zap.Error(err))
		return err
	}

	backgroundImage, err := s.lendingRepo.SelectBackgroundImageByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage s.lendingRepo.SelectBackgroundImageByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.lendingRepo.UpdateBackgroundImage(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage s.lendingRepo.UpdateBackgroundImage",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(backgroundImage)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + backgroundImage)
		if err != nil {
			s.logger.Error("internal.admin.UploadBackgroundImage os.Remove",
				zap.Error(err), zap.Int("landingID", landingID), zap.Any("old file", backgroundImage))
		}
	}

	return nil
}

func (s *service) UploadBackgroundForMobile(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile io.ReadAll", zap.Error(err))
		return err
	}

	backgroundImage, err := s.lendingRepo.SelectBackgroundForMobileByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile s.lendingRepo.SelectBackgroundForMobileByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.lendingRepo.UpdateBackgroundForMobile(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile s.lendingRepo.UpdateBackgroundForMobile",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(backgroundImage)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + backgroundImage)
		if err != nil {
			s.logger.Error("internal.admin.UploadBackgroundForMobile os.Remove",
				zap.Error(err), zap.Int("landingID", landingID), zap.Any("old file", backgroundImage))
		}
	}

	return nil
}

func (s *service) UploadMainLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo io.ReadAll", zap.Error(err))
		return err
	}

	oldFile, err := s.lendingRepo.SelectMainLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo s.lendingRepo.SelectMainLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.lendingRepo.UpdateMainLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo s.lendingRepo.UpdateMainLogo",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(oldFile)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + oldFile)
		if err != nil {
			s.logger.Error("internal.admin.UploadMainLogo os.Remove",
				zap.Error(err), zap.Int("landingID", landingID), zap.Any("old file", oldFile))
		}
	}

	return nil
}

func (s *service) UploadPartnerLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo io.ReadAll", zap.Error(err))
		return err
	}

	oldFile, err := s.lendingRepo.SelectPartnerLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo s.lendingRepo.SelectPartnerLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.lendingRepo.UpdatePartnerLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo s.lendingRepo.UpdatePartnerLogo",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(oldFile)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + oldFile)
		if err != nil {
			s.logger.Error("internal.admin.UploadPartnerLogo os.Remove",
				zap.Error(err), zap.Int("landingID", landingID), zap.Any("old file", oldFile))
		}
	}

	return nil
}

func (s *service) UploadOurLogo(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo io.ReadAll", zap.Error(err))
		return err
	}

	oldFile, err := s.lendingRepo.SelectOurLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo s.lendingRepo.SelectOurLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.lendingRepo.UpdateOurLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo s.lendingRepo.UpdateOurLogo",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(oldFile)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + oldFile)
		if err != nil {
			s.logger.Error("internal.admin.UploadOurLogo os.Remove",
				zap.Error(err), zap.Int("landingID", landingID), zap.Any("old file", oldFile))
		}
	}

	return nil
}

func (s *service) AddFeatureAndAmenity(ctx context.Context, file multipart.File, typeName, featureName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.AddFeatureAndAmenity os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.AddFeatureAndAmenity io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.AddFeatureAndAmenity io.ReadAll", zap.Error(err))
		return err
	}

	err = s.lendingRepo.InsertFeatureAndAmenity(ctx, featureName, filename)
	if err != nil {
		s.logger.Error("internal.admin.AddFeatureAndAmenity s.lendingRepo.UpdateOurLogo",
			zap.Error(err), zap.String("featureName", featureName))
		return err
	}

	return nil
}

func (s *service) DeleteFeatureAndAmenity(ctx context.Context, id int) error {

	list, err := s.lendingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, []int{id})
	if err != nil {
		s.logger.Error("internal.admin.DeleteFeatureAndAmenity s.lendingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	err = s.lendingRepo.DeleteFeatureAndAmenity(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.DeleteFeatureAndAmenity s.lendingRepo.DeleteFeatureAndAmenity",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	err = os.Remove(structs.FilePathRafflesHomes + list[0].Name)
	if err != nil {
		s.logger.Error("internal.admin.DeleteFeatureAndAmenity os.Remove",
			zap.Error(err), zap.Int("id", id), zap.Any("old file", list[0].Name))
	}

	return nil
}

func (s *service) GetFileURL(ctx context.Context) string {
	return structs.FileBaseURLPathRafflesHomes
}

func (s *service) GetSpecialGiftIcons(ctx context.Context) []structs.SpecialGiftIcon {
	return []structs.SpecialGiftIcon{
		{
			Gift: "Car",
			Icon: "car.svg",
		}, {
			Gift: "Rolex",
			Icon: "rolex.svg",
		},
	}
}
