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

func (s *service) GetLandingList(ctx context.Context, offset, limit int) (list []structs.LandingList, count int, err error) {
	list, count, err = s.landingRepo.GetLandingList(ctx, offset, limit)
	if err != nil {
		s.logger.Error("internal.admin.GetLandingList s.landingRepo.GetLandingList", zap.Error(err))
		return nil, 0, err
	}

	return list, 0, nil
}

func (s *service) GetLandingData(ctx context.Context, landingID int) (data structs.Landing, err error) {
	data, err = s.landingRepo.GetLandingByID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.GetLandingData s.landingRepo.GetLandingByID",
			zap.Error(err), zap.Int("landingID", landingID))
		return data, err
	}

	data.FeaturesAndAmenities, err = s.landingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, data.FeaturesAndAmenitiesIDs)
	if err != nil {
		s.logger.Error("internal.admin.GetLandingData s.landingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Int("landingID", landingID))
	}

	data.Availabilities, err = s.landingRepo.GetAvailabilitiesByLandingID(ctx, data.ID)
	if err != nil {
		s.logger.Error("internal.admin.GetLandingData s.landingRepo.GetAvailabilitiesByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
	}

	return data, nil
}

func (s *service) SaveLanding(ctx context.Context, data structs.Landing) error {
	err := s.landingRepo.SaveLanding(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.SaveLanding s.landingRepo.SaveLanding",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) UpdateLanding(ctx context.Context, data structs.Landing) error {
	err := s.landingRepo.UpdateLanding(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.UpdateLanding s.landingRepo.UpdateLanding",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) SaveAvailability(ctx context.Context, data structs.Availability) error {
	err := s.landingRepo.SaveAvailability(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.SaveAvailability s.landingRepo.SaveAvailability",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) UpdateAvailability(ctx context.Context, data structs.Availability) error {
	err := s.landingRepo.UpdateAvailability(ctx, data)
	if err != nil {
		s.logger.Error("internal.admin.UpdateAvailability s.landingRepo.UpdateAvailability",
			zap.Error(err), zap.Any("data", data))
		return err
	}

	return nil
}

func (s *service) RemoveAvailability(ctx context.Context, id int) error {
	err := s.landingRepo.DeleteAvailability(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.RemoveAvailability s.landingRepo.RemoveAvailability",
			zap.Error(err), zap.Any("id", id))
		return err
	}

	return nil
}

func (s *service) UploadLandingImages(ctx context.Context, id int, files []multipart.File) error {

	newImagesName, err := s.uploadImages(files, structs.FilePathRafflesHomes)
	if err != nil {
		s.logger.Error("internal.admin.UploadLandingImages s.uploadImages", zap.Error(err), zap.Int("id", id))
		return err
	}

	imagesName, err := s.landingRepo.GetImagesByLandingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.UploadLandingImages s.landingRepo.GetImagesByLandingID",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	imagesName = append(imagesName, newImagesName...)

	err = s.landingRepo.UpdateLandingImages(ctx, id, imagesName)
	if err != nil {
		s.logger.Error("internal.admin.UploadLandingImages s.landingRepo.UpdateLandingImages",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	return nil
}

func (s *service) DeleteLandingImages(ctx context.Context, id int, imageName string) error {
	imageNames, err := s.landingRepo.GetImagesByLandingID(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLandingImages s.landingRepo.GetImagesByLandingID",
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

	err = s.landingRepo.UpdateLandingImages(ctx, id, imageNames)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLandingImages s.landingRepo.UpdateLandingImages",
			zap.Error(err), zap.Int("id", id), zap.Any("images", imageNames))
		return err
	}

	err = os.Remove(structs.FilePathRafflesHomes + imageName)
	if err != nil {
		s.logger.Error("internal.admin.DeleteLandingImages os.Remove",
			zap.Error(err), zap.Int("id", id), zap.Any("image", imageName))
		return err
	}

	return nil
}

func (s *service) GetFeaturesAndAmenities(ctx context.Context) (list []structs.FeatureOrAmenity, err error) {
	list, err = s.landingRepo.SelectFeaturesAndAmenities(ctx)
	if err != nil {
		s.logger.Error("internal.admin.GetFeaturesAndAmenities s.landingRepo.SelectFeaturesAndAmenities",
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

	paymentPlan, err := s.landingRepo.SelectFilePlanByLandingID(ctx, availabilityID)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan s.landingRepo.SelectFilePlanByLandingID",
			zap.Error(err), zap.Int("availabilityID", availabilityID))
		return err
	}

	err = s.landingRepo.UpdateFilePlan(ctx, availabilityID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadFilePlan s.landingRepo.UpdateFilePlan",
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

	backgroundImage, err := s.landingRepo.SelectBackgroundImageByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage s.landingRepo.SelectBackgroundImageByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdateBackgroundImage(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundImage s.landingRepo.UpdateBackgroundImage",
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

	backgroundImage, err := s.landingRepo.SelectBackgroundForMobileByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile s.landingRepo.SelectBackgroundForMobileByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdateBackgroundForMobile(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadBackgroundForMobile s.landingRepo.UpdateBackgroundForMobile",
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

	oldFile, err := s.landingRepo.SelectMainLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo s.landingRepo.SelectMainLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdateMainLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadMainLogo s.landingRepo.UpdateMainLogo",
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

	oldFile, err := s.landingRepo.SelectPartnerLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo s.landingRepo.SelectPartnerLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdatePartnerLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadPartnerLogo s.landingRepo.UpdatePartnerLogo",
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

	oldFile, err := s.landingRepo.SelectOurLogoByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo s.landingRepo.SelectOurLogoByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdateOurLogo(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadOurLogo s.landingRepo.UpdateOurLogo",
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

func (s *service) UploadVideoCover(ctx context.Context, landingID int, file multipart.File, typeName string) error {
	filename := uuid.NewString() + typeName

	newFile, err := os.Create(structs.FilePathRafflesHomes + filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadVideoCover os.Create", zap.Error(err))
		return err
	}
	defer newFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error("internal.admin.UploadVideoCover io.ReadAll", zap.Error(err))
		return err
	}

	_, err = newFile.Write(fileBytes)
	if err != nil {
		s.logger.Error("internal.admin.UploadVideoCover io.ReadAll", zap.Error(err))
		return err
	}

	oldFile, err := s.landingRepo.SelectVideoCoverByLandingID(ctx, landingID)
	if err != nil {
		s.logger.Error("internal.admin.UploadVideoCover s.landingRepo.SelectVideoCoverByLandingID",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	err = s.landingRepo.UpdateVideoCover(ctx, landingID, filename)
	if err != nil {
		s.logger.Error("internal.admin.UploadVideoCover s.landingRepo.UpdateVideoCover",
			zap.Error(err), zap.Int("landingID", landingID))
		return err
	}

	if len(strings.TrimSpace(oldFile)) != 0 {
		err = os.Remove(structs.FilePathRafflesHomes + oldFile)
		if err != nil {
			s.logger.Error("internal.admin.UpdateVideoCover os.Remove",
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

	err = s.landingRepo.InsertFeatureAndAmenity(ctx, featureName, filename)
	if err != nil {
		s.logger.Error("internal.admin.AddFeatureAndAmenity s.landingRepo.UpdateOurLogo",
			zap.Error(err), zap.String("featureName", featureName))
		return err
	}

	return nil
}

func (s *service) DeleteFeatureAndAmenity(ctx context.Context, id int) error {

	list, err := s.landingRepo.SelectFeaturesAndAmenitiesByIDs(ctx, []int{id})
	if err != nil {
		s.logger.Error("internal.admin.DeleteFeatureAndAmenity s.landingRepo.SelectFeaturesAndAmenitiesByIDs",
			zap.Error(err), zap.Int("id", id))
		return err
	}

	err = s.landingRepo.DeleteFeatureAndAmenity(ctx, id)
	if err != nil {
		s.logger.Error("internal.admin.DeleteFeatureAndAmenity s.landingRepo.DeleteFeatureAndAmenity",
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
